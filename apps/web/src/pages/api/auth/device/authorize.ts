import type { APIRoute } from 'astro';
import * as auth from '../../../../lib/auth';

export const POST: APIRoute = async ({ request, locals, cookies }) => {
    const formData = await request.formData();
    const code = formData.get('code') as string;

    if (!code) {
        return new Response(null, {
            status: 302,
            headers: { 'Location': '/auth/device?error=invalid_code' }
        });
    }

    const token = cookies.get('auth_token')?.value;
    if (!token) {
        return new Response(null, {
            status: 302,
            headers: { 'Location': '/auth/login?next=/auth/device?code=' + code }
        });
    }

    const env = locals.runtime.env as auth.Env;
    const session = await auth.getSession(env, token);

    if (!session) {
        return new Response(null, {
            status: 302,
            headers: { 'Location': '/auth/login?next=/auth/device?code=' + code }
        });
    }

    const deviceCode = await auth.getDeviceCode(env, code);
    if (!deviceCode || deviceCode.status !== 'pending') {
        return new Response(null, {
            status: 302,
            headers: { 'Location': '/auth/device?error=invalid_code' }
        });
    }

    await auth.approveDeviceCode(env, code, session.userId, token);

    return new Response(null, {
        status: 302,
        headers: { 'Location': '/auth/device/success' }
    });
};
