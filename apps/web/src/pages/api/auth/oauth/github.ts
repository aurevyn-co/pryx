import type { APIRoute } from 'astro';

export const GET: APIRoute = async ({ url, locals, redirect }) => {
    const env = locals.runtime.env;
    const clientId = env.GITHUB_CLIENT_ID as string;
    
    const state = crypto.randomUUID();
    const host = url.host;
    const callbackUrl = `https://${host}/api/auth/oauth/github/callback`;
    const githubAuthUrl = new URL('https://github.com/login/oauth/authorize');
    githubAuthUrl.searchParams.set('client_id', clientId);
    githubAuthUrl.searchParams.set('redirect_uri', callbackUrl);
    githubAuthUrl.searchParams.set('scope', 'user:email');
    githubAuthUrl.searchParams.set('state', state);

    return new Response(null, {
        status: 302,
        headers: {
            'Location': githubAuthUrl.toString(),
            'Set-Cookie': `oauth_state=${state}; Path=/; Max-Age=600; HttpOnly; SameSite=Lax; Secure`
        }
    });
};
