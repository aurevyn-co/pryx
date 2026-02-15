import type { APIRoute } from 'astro';
import * as auth from '../../../../lib/auth';

export const POST: APIRoute = async ({ locals }) => {
    const env = locals.runtime.env as auth.Env;
    const { code, expiresIn } = await auth.createDeviceCode(env);

    return new Response(JSON.stringify({
        device_code: code,
        user_code: code,
        verification_uri: 'https://pryx.dev/auth/device',
        expires_in: expiresIn,
        interval: 5
    }), {
        headers: { 'Content-Type': 'application/json' }
    });
};
