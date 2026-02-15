import type { APIRoute } from 'astro';
import * as auth from '../../../lib/auth';

export const POST: APIRoute = async ({ request, locals, redirect }) => {
    const formData = await request.formData();
    const email = formData.get('email') as string;
    const password = formData.get('password') as string;
    const name = formData.get('name') as string;
    const confirm = formData.get('confirm') as string;
    const next = formData.get('next') as string || '/dashboard';

    if (!email || !password) {
        return redirect('/auth/register?error=invalid_email');
    }

    if (!auth.isValidEmail(email)) {
        return redirect('/auth/register?error=invalid_email');
    }

    const passwordValidation = auth.isValidPassword(password);
    if (!passwordValidation.valid) {
        return redirect('/auth/register?error=weak_password');
    }

    if (password !== confirm) {
        return redirect('/auth/register?error=password_mismatch');
    }

    const env = locals.runtime.env as auth.Env;
    const existingUser = await auth.getUserByEmail(env, email);

    if (existingUser) {
        return redirect('/auth/register?error=email_exists');
    }

    const user = await auth.createUser(env, email, password, name);
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
