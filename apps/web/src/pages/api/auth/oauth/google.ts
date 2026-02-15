import type { APIRoute } from 'astro';

export const GET: APIRoute = async ({ url, locals, redirect }) => {
    const env = locals.runtime.env as { GOOGLE_CLIENT_ID?: string };
    const clientId = env.GOOGLE_CLIENT_ID as string;
    
    const state = crypto.randomUUID();
    const host = url.host;
    const callbackUrl = `https://${host}/api/auth/oauth/google/callback`;
    
    const googleAuthUrl = new URL('https://accounts.google.com/o/oauth2/v2/auth');
    googleAuthUrl.searchParams.set('client_id', clientId);
    googleAuthUrl.searchParams.set('redirect_uri', callbackUrl);
    googleAuthUrl.searchParams.set('response_type', 'code');
    googleAuthUrl.searchParams.set('scope', 'openid email profile');
    googleAuthUrl.searchParams.set('state', state);
    googleAuthUrl.searchParams.set('access_type', 'offline');
    googleAuthUrl.searchParams.set('prompt', 'consent');

    return new Response(null, {
        status: 302,
        headers: {
            'Location': googleAuthUrl.toString(),
            'Set-Cookie': `oauth_state=${state}; Path=/; Max-Age=600; HttpOnly; SameSite=Lax; Secure`
        }
    });
};
