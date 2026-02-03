import { Hono } from 'hono';
import { handle } from 'hono/vercel';

// Direct Cloudflare Worker entry point
// This bypasses Astro's API route system and accesses bindings directly

const app = new Hono<{ Bindings: CloudflareBindings }>();

// Health check
app.get('/', (c) => {
    return c.json({ 
        name: 'Pryx Cloud API', 
        status: 'operational', 
        engine: 'hono',
        mode: 'direct-worker'
    });
});

// Device code endpoint
app.post('/api/auth/device/code', async (c) => {
    try {
        const formData = await c.req.formData();
        const deviceId = formData.get('device_id')?.toString() || '';
        const scopeStr = formData.get('scope')?.toString() || 'telemetry.write';

        // Generate codes
        const deviceCode = generateCode(40, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789');
        const userCode = `${generateCode(4, 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789')}-${generateCode(4, 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789')}`;

        const entry = {
            user_code: userCode,
            device_id: deviceId,
            scopes: scopeStr.split(' '),
            created_at: Date.now(),
            expires_at: Date.now() + 600 * 1000,
            authorized: false,
        };

        // Store in KV
        await c.env.DEVICE_CODES.put(deviceCode, JSON.stringify(entry), { expirationTtl: 660 });
        await c.env.DEVICE_CODES.put(`user:${userCode}`, deviceCode, { expirationTtl: 660 });

        return c.json({
            device_code: deviceCode,
            user_code: userCode,
            verification_uri: '/link',
            verification_uri_complete: `/link?code=${userCode}`,
            expires_in: 600,
            interval: 5,
        });
    } catch (error) {
        console.error('Device code error:', error);
        return c.json({ error: String(error) }, 500);
    }
});

// Helper function
function generateCode(length: number, charset: string): string {
    const array = new Uint8Array(length);
    crypto.getRandomValues(array);
    return Array.from(array, b => charset[b % charset.length]).join('');
}

// Export for Cloudflare Worker
export default {
    async fetch(request: Request, env: CloudflareBindings, ctx: ExecutionContext): Promise<Response> {
        return app.fetch(request, env, ctx);
    },
};

// Type definition for Cloudflare bindings
type CloudflareBindings = {
    DEVICE_CODES: KVNamespace;
    TOKENS: KVNamespace;
    SESSIONS: KVNamespace;
    RATE_LIMITER?: RateLimit;
    ENVIRONMENT?: string;
};
