import type { APIRoute } from 'astro';
import * as auth from '../../../../lib/auth';

export const POST: APIRoute = async ({ request, locals, redirect }) => {
    const formData = await request.formData();
    const token = formData.get('token') as string;
    const password = formData.get('password') as string;
    const confirm = formData.get('confirm') as string;

    if (!token || !password || password !== confirm) {
        return redirect('/auth/reset?error=invalid_request');
    }

    const passwordValidation = auth.isValidPassword(password);
    if (!passwordValidation.valid) {
        return redirect('/auth/reset?error=weak_password');
    }

    const env = locals.runtime.env as auth.Env;
    const userId = await auth.verifyResetToken(env, token);

    if (!userId) {
        return redirect('/auth/reset?error=invalid_token');
    }

    await auth.updatePassword(env, userId, password);
    await auth.useResetToken(env, token);

    return redirect('/auth/login?success=password_reset');
};
