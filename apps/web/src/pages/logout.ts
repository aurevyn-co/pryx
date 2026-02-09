import type { APIRoute } from 'astro';

export const GET: APIRoute = async ({ url }) => {
    const next = url.searchParams.get('next') || '/auth';

    const headers = new Headers({
        Location: next,
    });

    headers.append('Set-Cookie', 'auth_token=; Path=/; Max-Age=0; SameSite=Lax');
    headers.append('Set-Cookie', 'auth_role=; Path=/; Max-Age=0; SameSite=Lax');

    return new Response(null, {
        status: 302,
        headers,
    });
};
