// Admin API routes for superadmin dashboard
// These endpoints provide global telemetry and fleet management data

import { Hono } from 'hono';
import { cors } from 'hono/cors';

// ============================================================================
// 3-LAYER ARCHITECTURE TYPES
// ============================================================================

/**
 * Layer types for role-based access control
 * - user: Regular cloud users viewing their own data
 * - superadmin: Full admin access to all devices/users/telemetry
 * - localhost: Local admin via runtime port (no auth required)
 */
type Layer = 'user' | 'superadmin' | 'localhost';

/**
 * Layer context containing authentication and authorization information
 */
interface LayerContext {
  layer: Layer;
  userId?: string;
  isLocalhost: boolean;
}

/**
 * Extended context with layer information attached
 */
interface AdminContext extends LayerContext {
  env: any;
}

// ============================================================================
// LAYER DETECTION MIDDLEWARE
// ============================================================================

/**
 * Extract layer information from Authorization header
 * Pattern: Bearer {layer}:{identifier}
 * - Bearer user:user-123 → User layer, user ID user-123
 * - Bearer superadmin:admin-key → Superadmin layer
 * - Bearer localhost or no header → Localhost layer
 */
function extractLayer(authHeader: string | null, env: Record<string, string | undefined>): LayerContext {
  // Check for localhost bypass via environment variable
  const localhostKey = env.LOCALHOST_ADMIN_KEY || process.env.LOCALHOST_ADMIN_KEY;

  if (!authHeader) {
    // No auth header means localhost layer (local control panel)
    return { layer: 'localhost', isLocalhost: true };
  }

  const token = authHeader.replace('Bearer ', '').trim();

  // Check for localhost bypass key
  if (token === localhostKey || token === 'localhost') {
    return { layer: 'localhost', isLocalhost: true };
  }

  // Check for superadmin token
  if (token.startsWith('superadmin:')) {
    const adminId = token.replace('superadmin:', '');
    return { layer: 'superadmin', userId: adminId, isLocalhost: false };
  }

  // Check for regular user token
  if (token.startsWith('user:')) {
    const userId = token.replace('user:', '');
    return { layer: 'user', userId, isLocalhost: false };
  }

  // Fallback: treat as user with token as userId
  return { layer: 'user', userId: token, isLocalhost: false };
}

/**
 * Require a specific layer or higher for access
 * @param requiredLayers - Layers that are allowed access
 */
function requireLayer(...requiredLayers: Layer[]) {
  return async (c: any, next: () => Promise<void>) => {
    const authHeader = c.req.header('Authorization');
    const env = c.env as Record<string, string | undefined>;
    const layerContext = extractLayer(authHeader, env);

    // Store layer context in request for downstream handlers
    (c as any).req.layerContext = layerContext;

    // Check if user's layer is allowed
    if (!requiredLayers.includes(layerContext.layer)) {
      const layerNames: Record<Layer, string> = {
        user: 'user',
        superadmin: 'superadmin',
        localhost: 'localhost',
      };

      const allowedStr = requiredLayers.map(l => layerNames[l]).join(' or ');
      const actualStr = layerNames[layerContext.layer];

      return c.json({
        error: 'Forbidden',
        message: `This endpoint requires ${allowedStr} access. You have ${actualStr} access.`,
        required: requiredLayers,
        current: layerContext.layer,
      }, 403);
    }

    await next();
  };
}

/**
 * Get layer-aware filtering for user-scoped data
 * Returns filter criteria based on the user's layer
 */
function getLayerFilters(layerContext: LayerContext): { userId?: string; global?: boolean } {
  switch (layerContext.layer) {
    case 'user':
      // Regular users can only see their own data
      return { userId: layerContext.userId };
    case 'superadmin':
      // Superadmins can see everything
      return { global: true };
    case 'localhost':
      // Localhost can see everything (local control panel)
      return { global: true };
    default:
      return { global: true };
  }
}

// Admin API router
const adminApi = new Hono<{ Bindings: any }>();

// Enable CORS for admin API
adminApi.use('/*', cors({
  origin: ['http://localhost:4321', 'https://pryx.dev'],
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE'],
  allowHeaders: ['Content-Type', 'Authorization'],
}));

// Middleware to verify admin authentication
adminApi.use('/*', async (c, next) => {
  const authHeader = c.req.header('Authorization');

  // TODO: Implement proper admin authentication
  // For now, check for admin API key
  const env = c.env as Record<string, string | undefined>;
  const adminKey = env.ADMIN_API_KEY || process.env.ADMIN_API_KEY;

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return c.json({ error: 'Unauthorized - Missing bearer token' }, 401);
  }

  const token = authHeader.replace('Bearer ', '');

  // Validate admin token
  if (token !== adminKey) {
    return c.json({ error: 'Unauthorized - Invalid credentials' }, 401);
  }

  await next();
});

adminApi.get('/stats', requireLayer('superadmin', 'localhost'), async (c) => {
  const layerContext = (c as any).req.layerContext as LayerContext;
  const filters = getLayerFilters(layerContext);

  let stats: Record<string, any>;

  try {
    if (filters.global) {
      const list = await c.env.TELEMETRY.list({ limit: 1000, prefix: 'telemetry:' });

      let totalEvents = list.keys.length;
      let totalCost = 0;
      let errorCount = 0;
      const uniqueDevices = new Set<string>();
      const uniqueSessions = new Set<string>();

      for (const key of list.keys) {
        try {
          const value = await c.env.TELEMETRY.get(key.name);
          if (value) {
            const event = JSON.parse(value);

            if (event.device_id) uniqueDevices.add(event.device_id);
            if (event.session_id) uniqueSessions.add(event.session_id);
            if (event.cost) totalCost += event.cost;
            if (event.level === 'error' || event.level === 'critical') errorCount++;
          }
        } catch (e) {
        }
      }

      const now = Date.now();
      const oneDayAgo = now - 86400000;
      const recentEvents = await c.env.TELEMETRY.list({ limit: 100, prefix: 'telemetry:' });

      let newEventsToday = 0;
      for (const key of recentEvents.keys) {
        try {
          const value = await c.env.TELEMETRY.get(key.name);
          if (value) {
            const event = JSON.parse(value);
            if (event.received_at && event.received_at >= oneDayAgo) {
              newEventsToday++;
            }
          }
        } catch (e) {
        }
      }

      stats = {
        totalEvents,
        uniqueDevices: uniqueDevices.size,
        uniqueSessions: uniqueSessions.size,
        totalCost,
        errorCount,
        errorRate: totalEvents > 0 ? errorCount / totalEvents : 0,
        newEventsToday,
        timestamp: new Date().toISOString(),
      };
    } else {
      stats = {
        message: 'User-specific stats require user registry implementation',
        timestamp: new Date().toISOString(),
      };
    }

    return c.json(stats);
  } catch (e) {
    console.error('Stats error:', e);
    return c.json({ error: String(e) }, 500);
  }
});

// GET /api/admin/users - List users (layer-aware)
adminApi.get('/users', requireLayer('superadmin', 'localhost', 'user'), async (c) => {
  const layerContext = (c as any).req.layerContext as LayerContext;
  const range = c.req.query('range') || '7d';
  const filters = getLayerFilters(layerContext);

  // Layer-aware user listing
  let users: Array<Record<string, any>>;

  if (filters.userId && layerContext.layer === 'user') {
    // Regular user: return only their own user data
    users = [
      {
        id: filters.userId,
        email: 'user@example.com',
        createdAt: '2026-01-15T10:30:00Z',
        lastActive: '2026-02-03T14:22:00Z',
        deviceCount: 3,
        sessionCount: 156,
        totalCost: 45.20,
        status: 'active',
      },
    ];
  } else {
    // Superadmin or localhost: return all users
    users = [
      {
        id: 'user-001',
        email: 'admin@pryx.dev',
        createdAt: '2026-01-15T10:30:00Z',
        lastActive: '2026-02-03T14:22:00Z',
        deviceCount: 3,
        sessionCount: 156,
        totalCost: 45.20,
        status: 'active',
      },
      {
        id: 'user-002',
        email: 'demo@example.com',
        createdAt: '2026-01-20T08:15:00Z',
        lastActive: '2026-02-03T09:45:00Z',
        deviceCount: 2,
        sessionCount: 89,
        totalCost: 12.50,
        status: 'active',
      },
    ];
  }

  return c.json(users);
});

// GET /api/admin/users/:id - Get detailed user info
adminApi.get('/users/:id', async (c) => {
  const userId = c.req.param('id');

  // TODO: Fetch from D1 database
  const user = {
    id: userId,
    email: 'user@example.com',
    createdAt: '2026-01-15T10:30:00Z',
    lastActive: '2026-02-03T14:22:00Z',
    deviceCount: 3,
    sessionCount: 156,
    totalCost: 45.20,
    status: 'active',
    devices: [
      {
        id: 'dev-001',
        name: 'MacBook Pro',
        platform: 'macos',
        version: '1.0.0',
        status: 'online',
      },
    ],
    providers: ['openai', 'anthropic'],
    channels: ['telegram'],
  };

  return c.json(user);
});

// PUT /api/admin/users/:id - Update user (suspend/activate)
adminApi.put('/users/:id', async (c) => {
  const userId = c.req.param('id');
  const body = await c.req.json();

  // TODO: Update user status in database
  return c.json({
    id: userId,
    status: body.status,
    updatedAt: new Date().toISOString(),
  });
});

// GET /api/admin/devices - List devices (layer-aware)
adminApi.get('/devices', requireLayer('superadmin', 'localhost', 'user'), async (c) => {
  const layerContext = (c as any).req.layerContext as LayerContext;
  const filters = getLayerFilters(layerContext);

  // Layer-aware device listing
  let devices: Array<Record<string, any>>;

  if (filters.userId && layerContext.layer === 'user') {
    // Regular user: return only their own devices
    devices = [
      {
        id: 'dev-001',
        userId: filters.userId,
        userEmail: 'user@example.com',
        name: 'MacBook Pro',
        platform: 'macos',
        version: '1.0.0',
        status: 'online',
        lastSeen: '2026-02-03T14:22:00Z',
        ipAddress: '192.168.1.100',
      },
    ];
  } else {
    // Superadmin or localhost: return all devices
    devices = [
      {
        id: 'dev-001',
        userId: 'user-001',
        userEmail: 'admin@pryx.dev',
        name: 'MacBook Pro',
        platform: 'macos',
        version: '1.0.0',
        status: 'online',
        lastSeen: '2026-02-03T14:22:00Z',
        ipAddress: '192.168.1.100',
      },
      {
        id: 'dev-002',
        userId: 'user-001',
        userEmail: 'admin@pryx.dev',
        name: 'iPhone 15',
        platform: 'ios',
        version: '1.0.0',
        status: 'offline',
        lastSeen: '2026-02-03T10:15:00Z',
        ipAddress: null,
      },
    ];
  }

  return c.json(devices);
});

// POST /api/admin/devices/:id/sync - Force sync a device
adminApi.post('/devices/:id/sync', async (c) => {
  const deviceId = c.req.param('id');

  // TODO: Trigger device sync
  return c.json({
    id: deviceId,
    syncStatus: 'initiated',
    timestamp: new Date().toISOString(),
  });
});

// POST /api/admin/devices/:id/unpair - Unpair a device
adminApi.post('/devices/:id/unpair', async (c) => {
  const deviceId = c.req.param('id');

  // TODO: Unpair device in database
  return c.json({
    id: deviceId,
    status: 'unpaired',
    timestamp: new Date().toISOString(),
  });
});

// GET /api/admin/costs - Cost analytics (layer-aware)
adminApi.get('/costs', requireLayer('superadmin', 'localhost', 'user'), async (c) => {
  const layerContext = (c as any).req.layerContext as LayerContext;
  const range = c.req.query('range') || '7d';
  const filters = getLayerFilters(layerContext);

  // Layer-aware cost analytics
  let costs: Record<string, any>;

  if (filters.userId && layerContext.layer === 'user') {
    // Regular user: return their own cost data
    costs = {
      total: 45.20,
      byProvider: {
        openai: 25.00,
        anthropic: 20.20,
      },
      byDay: [
        { date: '2026-02-01', cost: 12.50 },
        { date: '2026-02-02', cost: 8.20 },
        { date: '2026-02-03', cost: 24.50 },
      ],
    };
  } else {
    // Superadmin or localhost: return global cost data
    costs = {
      total: 2847.50,
      byProvider: {
        openai: 1250.00,
        anthropic: 980.50,
        google: 617.00,
      },
      byDay: [
        { date: '2026-02-01', cost: 450.00 },
        { date: '2026-02-02', cost: 380.50 },
        { date: '2026-02-03', cost: 210.00 },
      ],
      topUsers: [
        { userId: 'user-001', email: 'admin@pryx.dev', cost: 45.20 },
        { userId: 'user-002', email: 'demo@example.com', cost: 12.50 },
      ],
    };
  }

  return c.json(costs);
});

adminApi.get('/health', async (c) => {
  const startTime = Date.now();
  let dbStatus = 'connected';
  let errorRate = 0;
  let telemetryCount = 0;

  try {
    const list = await c.env.TELEMETRY.list({ limit: 1 });
    const listTime = Date.now() - startTime;

    const recentTelemetry = await c.env.TELEMETRY.list({
      limit: 100,
      prefix: 'telemetry:',
    });

    telemetryCount = recentTelemetry.keys.length;

    let errorCount = 0;
    for (const key of recentTelemetry.keys.slice(0, 20)) {
      try {
        const value = await c.env.TELEMETRY.get(key.name);
        if (value) {
          const event = JSON.parse(value);
          if (event.level === 'error' || event.level === 'critical') {
            errorCount++;
          }
        }
      } catch (e) {
      }
    }
    errorRate = errorCount / Math.min(20, telemetryCount) || 0;

    const health = {
      runtimeStatus: errorRate > 0.1 ? 'degraded' : 'healthy',
      apiLatency: listTime,
      errorRate,
      dbStatus,
      queueDepth: 0,
      activeConnections: telemetryCount,
      timestamp: new Date().toISOString(),
    };

    return c.json(health);
  } catch (e) {
    console.error('Health check error:', e);
    const health = {
      runtimeStatus: 'critical',
      apiLatency: 9999,
      errorRate: 1.0,
      dbStatus: 'disconnected',
      queueDepth: 0,
      activeConnections: 0,
      timestamp: new Date().toISOString(),
    };

    return c.json(health);
  }
});

adminApi.get('/telemetry', async (c) => {
  const limit = parseInt(c.req.query('limit') || '50');
  const level = c.req.query('level');

  try {
    const list = await c.env.TELEMETRY.list({
      limit,
      prefix: 'telemetry:',
    });

    const events = [];
    for (const key of list.keys) {
      try {
        const value = await c.env.TELEMETRY.get(key.name);
        if (value) {
          const event = JSON.parse(value);
          if (!level || event.level === level) {
            events.push(event);
          }
        }
      } catch (e) {
        console.error('Failed to parse telemetry event:', e);
      }
    }

    return c.json({
      count: events.length,
      events: events.slice(0, limit),
    });
  } catch (e) {
    console.error('Telemetry query error:', e);
    return c.json({ error: String(e) }, 500);
  }
});

// GET /api/admin/telemetry/config - Get telemetry configuration (superadmin only)
adminApi.get('/telemetry/config', requireLayer('superadmin', 'localhost'), async (c) => {
  const config = {
    enabled: true,
    retentionDays: 7,
    exportBackend: 'https://api.pryx.dev',
    samplingRate: 0.1,
    logLevel: 'info',
    sensitiveDataFiltering: true,
    batchSize: 100,
    flushInterval: 5000,
  };

  return c.json(config);
});

// PUT /api/admin/telemetry/config - Update telemetry configuration (superadmin only)
adminApi.put('/telemetry/config', requireLayer('superadmin', 'localhost'), async (c) => {
  const body = await c.req.json();

  const config = {
    enabled: body.enabled ?? true,
    retentionDays: body.retentionDays ?? 7,
    exportBackend: body.exportBackend ?? 'https://api.pryx.dev',
    samplingRate: body.samplingRate ?? 0.1,
    logLevel: body.logLevel ?? 'info',
    sensitiveDataFiltering: body.sensitiveDataFiltering ?? true,
    batchSize: body.batchSize ?? 100,
    flushInterval: body.flushInterval ?? 5000,
    updatedAt: new Date().toISOString(),
    updatedBy: (c as any).req.layerContext?.userId || 'localhost',
  };

  return c.json(config);
});

adminApi.get('/logs', async (c) => {
  const level = c.req.query('level') || 'info';
  const limit = parseInt(c.req.query('limit') || '100');

  try {
    const list = await c.env.TELEMETRY.list({
      limit,
      prefix: 'telemetry:',
    });

    const logs = [];
    for (const key of list.keys) {
      try {
        const value = await c.env.TELEMETRY.get(key.name);
        if (value) {
          const event = JSON.parse(value);
          if (!level || event.level === level) {
            logs.push({
              timestamp: event.received_at ? new Date(event.received_at).toISOString() : new Date().toISOString(),
              level: event.level || 'info',
              message: event.message || event.type || 'Telemetry event',
              userId: event.user_id,
              deviceId: event.device_id,
              error: event.error,
            });
          }
        }
      } catch (e) {
        console.error('Failed to parse log entry:', e);
      }
    }

    return c.json({ logs: logs.slice(0, limit), total: logs.length });
  } catch (e) {
    console.error('Logs query error:', e);
    return c.json({ error: String(e) }, 500);
  }
});

export default adminApi;
