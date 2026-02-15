export interface Env {
    DB: D1Database;
    SESSION: KVNamespace;
    DEVICE_CODES: KVNamespace;
}

export interface User {
    id: string;
    email: string;
    name: string | null;
    email_verified: boolean;
    created_at: string;
}

export interface Session {
    id: string;
    user_id: string;
    token_hash: string;
    expires_at: string;
}

export function generateId(): string {
    return crypto.randomUUID();
}

export async function hashPassword(password: string): Promise<string> {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = new Uint8Array(hashBuffer);
    return Array.from(hashArray, b => b.toString(16).padStart(2, '0')).join('');
}

export async function verifyPassword(password: string, hash: string): Promise<boolean> {
    const passwordHash = await hashPassword(password);
    return passwordHash === hash;
}

export function generateToken(): string {
    const bytes = new Uint8Array(32);
    crypto.getRandomValues(bytes);
    return Array.from(bytes, b => b.toString(16).padStart(2, '0')).join('');
}

export function hashToken(token: string): string {
    const encoder = new TextEncoder();
    const data = encoder.encode(token);
    let hash = 0;
    for (let i = 0; i < data.length; i++) {
        const char = data[i];
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash;
    }
    return Math.abs(hash).toString(16).padStart(16, '0');
}

export async function createSession(env: Env, userId: string): Promise<string> {
    const token = generateToken();
    const tokenHash = hashToken(token);
    const sessionId = generateId();
    const expiresAt = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString();

    await env.DB.prepare(
        'INSERT INTO sessions (id, user_id, token_hash, expires_at) VALUES (?, ?, ?, ?)'
    ).bind(sessionId, userId, tokenHash, expiresAt).run();

    await env.SESSION.put(tokenHash, JSON.stringify({ userId, sessionId }), {
        expirationTtl: 30 * 24 * 60 * 60
    });

    return token;
}

export async function getSession(env: Env, token: string): Promise<{ userId: string; sessionId: string } | null> {
    const tokenHash = hashToken(token);
    const cached = await env.SESSION.get(tokenHash);
    
    if (cached) {
        return JSON.parse(cached);
    }

    const result = await env.DB.prepare(
        'SELECT user_id, id as session_id FROM sessions WHERE token_hash = ? AND expires_at > datetime("now")'
    ).bind(tokenHash).first();

    if (result) {
        await env.SESSION.put(tokenHash, JSON.stringify({ userId: result.user_id as string, sessionId: result.id as string }), {
            expirationTtl: 30 * 24 * 60 * 60
        });
        return { userId: result.user_id as string, sessionId: result.id as string };
    }

    return null;
}

export async function deleteSession(env: Env, token: string): Promise<void> {
    const tokenHash = hashToken(token);
    await env.DB.prepare('DELETE FROM sessions WHERE token_hash = ?').bind(tokenHash).run();
    await env.SESSION.delete(tokenHash);
}

export async function getUserById(env: Env, userId: string): Promise<User | null> {
    const result = await env.DB.prepare(
        'SELECT id, email, name, email_verified, created_at FROM users WHERE id = ?'
    ).bind(userId).first();

    if (!result) return null;

    return {
        id: result.id as string,
        email: result.email as string,
        name: result.name as string | null,
        email_verified: !!result.email_verified,
        created_at: result.created_at as string
    };
}

export async function getUserByEmail(env: Env, email: string): Promise<User | null> {
    const result = await env.DB.prepare(
        'SELECT id, email, name, email_verified, created_at FROM users WHERE email = ?'
    ).bind(email.toLowerCase()).first();

    if (!result) return null;

    return {
        id: result.id as string,
        email: result.email as string,
        name: result.name as string | null,
        email_verified: !!result.email_verified,
        created_at: result.created_at as string
    };
}

export async function createUser(env: Env, email: string, password: string, name?: string): Promise<User> {
    const id = generateId();
    const passwordHash = await hashPassword(password);

    await env.DB.prepare(
        'INSERT INTO users (id, email, password_hash, name) VALUES (?, ?, ?, ?)'
    ).bind(id, email.toLowerCase(), passwordHash, name || null).run();

    return {
        id,
        email: email.toLowerCase(),
        name: name || null,
        email_verified: false,
        created_at: new Date().toISOString()
    };
}

export async function verifyUserPassword(env: Env, email: string, password: string): Promise<User | null> {
    const result = await env.DB.prepare(
        'SELECT id, email, password_hash, name, email_verified, created_at FROM users WHERE email = ?'
    ).bind(email.toLowerCase()).first();

    if (!result) return null;

    const isValid = await verifyPassword(password, result.password_hash as string);
    if (!isValid) return null;

    return {
        id: result.id as string,
        email: result.email as string,
        name: result.name as string | null,
        email_verified: !!result.email_verified,
        created_at: result.created_at as string
    };
}

export async function createResetToken(env: Env, userId: string): Promise<string> {
    const token = generateToken();
    const tokenHash = hashToken(token);
    const id = generateId();
    const expiresAt = new Date(Date.now() + 60 * 60 * 1000).toISOString();

    await env.DB.prepare(
        'INSERT INTO reset_tokens (id, user_id, token_hash, expires_at) VALUES (?, ?, ?, ?)'
    ).bind(id, userId, tokenHash, expiresAt).run();

    return token;
}

export async function verifyResetToken(env: Env, token: string): Promise<string | null> {
    const tokenHash = hashToken(token);
    const result = await env.DB.prepare(
        'SELECT user_id FROM reset_tokens WHERE token_hash = ? AND expires_at > datetime("now") AND used = 0'
    ).bind(tokenHash).first();

    return result?.user_id as string | null;
}

export async function useResetToken(env: Env, token: string): Promise<void> {
    const tokenHash = hashToken(token);
    await env.DB.prepare('UPDATE reset_tokens SET used = 1 WHERE token_hash = ?').bind(tokenHash).run();
}

export async function updatePassword(env: Env, userId: string, password: string): Promise<void> {
    const passwordHash = await hashPassword(password);
    await env.DB.prepare('UPDATE users SET password_hash = ?, updated_at = datetime("now") WHERE id = ?')
        .bind(passwordHash, userId).run();
}

export async function createDeviceCode(env: Env): Promise<{ code: string; expiresIn: number }> {
    const code = Math.random().toString(36).substring(2, 8).toUpperCase();
    const expiresIn = 600;

    await env.DEVICE_CODES.put(code, JSON.stringify({ status: 'pending' }), {
        expirationTtl: expiresIn
    });

    return { code, expiresIn };
}

export async function getDeviceCode(env: Env, code: string): Promise<{ status: string; userId?: string; token?: string } | null> {
    const data = await env.DEVICE_CODES.get(code);
    return data ? JSON.parse(data) : null;
}

export async function approveDeviceCode(env: Env, code: string, userId: string, token: string): Promise<void> {
    await env.DEVICE_CODES.put(code, JSON.stringify({ status: 'approved', userId, token }), {
        expirationTtl: 600
    });
}

export function isValidEmail(email: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

export function isValidPassword(password: string): { valid: boolean; error?: string } {
    if (password.length < 8) {
        return { valid: false, error: 'Password must be at least 8 characters' };
    }
    if (!/[A-Z]/.test(password)) {
        return { valid: false, error: 'Password must contain an uppercase letter' };
    }
    if (!/[a-z]/.test(password)) {
        return { valid: false, error: 'Password must contain a lowercase letter' };
    }
    if (!/[0-9]/.test(password)) {
        return { valid: false, error: 'Password must contain a number' };
    }
    return { valid: true };
}
