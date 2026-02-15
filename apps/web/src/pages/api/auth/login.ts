import type { APIRoute } from 'astro';
import * as auth from '../../../lib/auth';

export const POST: APIRoute = async ({ request, locals, redirect }) => {
    const formData = await request.formData();
    const email = formData.get('email') as string;
    const password = formData.get('password') as string;
    const next = formData.get('next') as string || '/dashboard';

    if (!email || !password) {
        return redirect('/auth/login?error=invalid_credentials');
    }

    const env = locals.runtime.env as auth.Env;
    const user = await auth.verifyUserPassword(env, email, password);

    if (!user) {
        return redirect('/auth/login?error=invalid_credentials');
    }

    const token = await auth.createSession(env, user.id);
    const isSecure = request.url.startsWith('https');
    
    return new Response(null, {
        status: 302,
        headers: {
            'Location': next,
            'Set-Cookie': `auth_token=${token}; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax${isSecure ? '; Secure' : ''}`
        }
    });
};
