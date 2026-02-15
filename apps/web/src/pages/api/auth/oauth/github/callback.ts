import type { APIRoute } from 'astro';
import * as auth from '../../../../../lib/auth';

export const GET: APIRoute = async ({ url, locals, cookies, redirect }) => {
    try {
        const env = locals.runtime.env as auth.Env & { GITHUB_CLIENT_ID: string; GITHUB_CLIENT_SECRET: string };
        
        const code = url.searchParams.get('code');
        const state = url.searchParams.get('state');
        const storedState = cookies.get('oauth_state')?.value;

        if (!code || !state || state !== storedState) {
            return redirect('/auth/login?error=oauth_failed');
        }

        const clientId = env.GITHUB_CLIENT_ID;
        const clientSecret = env.GITHUB_CLIENT_SECRET;

        const tokenResponse = await fetch('https://github.com/login/oauth/access_token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                client_id: clientId,
                client_secret: clientSecret,
                code
            })
        });

        const tokenData = await tokenResponse.json() as { access_token?: string; error?: string; error_description?: string };
        
        if (!tokenData.access_token) {
            console.error('GitHub OAuth error:', tokenData);
            return redirect('/auth/login?error=oauth_failed');
        }

        const userResponse = await fetch('https://api.github.com/user', {
            headers: {
                'Authorization': `token ${tokenData.access_token}`,
                'Accept': 'application/vnd.github.v3+json'
            }
        });

        const githubUser = await userResponse.json() as { id: number; email?: string; name?: string; login: string };

        let email = githubUser.email;
        if (!email) {
            const emailsResponse = await fetch('https://api.github.com/user/emails', {
                headers: {
                    'Authorization': `token ${tokenData.access_token}`,
                    'Accept': 'application/vnd.github.v3+json'
                }
            });
            const emails = await emailsResponse.json() as Array<{ email: string; primary: boolean; verified: boolean }>;
            const primaryEmail = emails.find(e => e.primary && e.verified);
            email = primaryEmail?.email;
        }

        if (!email) {
            return redirect('/auth/login?error=no_email');
        }

        const existingOAuth = await env.DB.prepare(
            'SELECT user_id FROM oauth_accounts WHERE provider = ? AND provider_user_id = ?'
        ).bind('github', String(githubUser.id)).first();

        let userId: string;

        if (existingOAuth) {
            userId = existingOAuth.user_id as string;
        } else {
            const existingUser = await auth.getUserByEmail(env, email);
            
            if (existingUser) {
                await env.DB.prepare(
                    'INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id) VALUES (?, ?, ?, ?)'
                ).bind(crypto.randomUUID(), existingUser.id, 'github', String(githubUser.id)).run();
                userId = existingUser.id;
            } else {
                const newUser = await auth.createUser(env, email, crypto.randomUUID(), githubUser.name || githubUser.login);
                await env.DB.prepare(
                    'INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id) VALUES (?, ?, ?, ?)'
                ).bind(crypto.randomUUID(), newUser.id, 'github', String(githubUser.id)).run();
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
