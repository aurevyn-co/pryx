import type { APIRoute } from 'astro';
import * as auth from '../../../../../lib/auth';

export const GET: APIRoute = async ({ url, locals, cookies, redirect }) => {
    try {
        const env = locals.runtime.env as auth.Env & { GOOGLE_CLIENT_ID: string; GOOGLE_CLIENT_SECRET: string };
        
        const code = url.searchParams.get('code');
        const state = url.searchParams.get('state');
        const storedState = cookies.get('oauth_state')?.value;

        if (!code || !state || state !== storedState) {
            return redirect('/auth/login?error=oauth_failed');
        }

        const clientId = env.GOOGLE_CLIENT_ID;
        const clientSecret = env.GOOGLE_CLIENT_SECRET;

        const tokenResponse = await fetch('https://oauth2.googleapis.com/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: new URLSearchParams({
                client_id: clientId,
                client_secret: clientSecret,
                code,
                grant_type: 'authorization_code',
                redirect_uri: 'https://pryx.dev/api/auth/oauth/google/callback'
            })
        });

        const tokenData = await tokenResponse.json() as { 
            access_token?: string; 
            id_token?: string;
            error?: string; 
            error_description?: string 
        };
        
        if (!tokenData.access_token) {
            console.error('Google OAuth error:', tokenData);
            return redirect('/auth/login?error=oauth_failed');
        }

        // Decode the ID token to get user info
        const idTokenParts = tokenData.id_token?.split('.') || [];
        let userInfo: { email?: string; name?: string; picture?: string; sub?: string } = {};
        
        if (idTokenParts.length >= 2) {
            try {
                const payload = JSON.parse(atob(idTokenParts[1]));
                userInfo = {
                    email: payload.email,
                    name: payload.name,
                    picture: payload.picture,
                    sub: payload.sub
                };
            } catch {
                console.error('Failed to decode ID token');
            }
        }

        if (!userInfo.email) {
            return redirect('/auth/login?error=no_email');
        }

        const existingOAuth = await env.DB.prepare(
            'SELECT user_id FROM oauth_accounts WHERE provider = ? AND provider_user_id = ?'
        ).bind('google', userInfo.sub).first();

        let userId: string;

        if (existingOAuth) {
            userId = existingOAuth.user_id as string;
        } else {
            const existingUser = await auth.getUserByEmail(env, userInfo.email);
            
            if (existingUser) {
                await env.DB.prepare(
                    'INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id) VALUES (?, ?, ?, ?)'
                ).bind(crypto.randomUUID(), existingUser.id, 'google', userInfo.sub).run();
                userId = existingUser.id;
            } else {
                const newUser = await auth.createUser(env, userInfo.email, crypto.randomUUID(), userInfo.name || undefined);
                await env.DB.prepare(
                    'INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id) VALUES (?, ?, ?, ?)'
                ).bind(crypto.randomUUID(), newUser.id, 'google', userInfo.sub).run();
                userId = newUser.id;
            }
        }

        const token = await auth.createSession(env, userId);
        
        return new Response(null, {
            status: 302,
            headers: {
                'Location': '/dashboard',
                'Set-Cookie': `auth_token=${token}; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax; Secure`
            }
        });
    } catch (err) {
        console.error('OAuth callback error:', err);
        return redirect('/auth/login?error=server_error');
    }
};
