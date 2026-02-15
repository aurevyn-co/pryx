import type { APIRoute } from 'astro';
import * as auth from '../../../lib/auth';

export const POST: APIRoute = async ({ request, locals, redirect }) => {
    const formData = await request.formData();
    const email = formData.get('email') as string;

    if (!email || !auth.isValidEmail(email)) {
        return redirect('/auth/reset?error=invalid_email');
    }

    const env = locals.runtime.env as auth.Env;
    const user = await auth.getUserByEmail(env, email);

    if (!user) {
        return redirect('/auth/reset?error=not_found');
    }

    const token = await auth.createResetToken(env, user.id);
    console.log(`Password reset token for ${email}: ${token}`);

    return redirect('/auth/reset?success=1');
};

export const GET: APIRoute = async ({ url, locals, redirect }) => {
    const token = url.searchParams.get('token');
    
    if (!token) {
        return redirect('/auth/reset?error=invalid_token');
    }

    const env = locals.runtime.env as auth.Env;
    const userId = await auth.verifyResetToken(env, token);

    if (!userId) {
        return redirect('/auth/reset?error=invalid_token');
    }

    return new Response(`
        <!DOCTYPE html>
        <html>
        <head><title>Reset Password</title></head>
        <body style="background:#0f1730;color:#eef3ff;font-family:system-ui;display:flex;align-items:center;justify-content:center;min-height:100vh">
            <form action="/api/auth/reset/confirm" method="POST" style="background:rgba(9,16,30,0.9);padding:2rem;border-radius:1rem;border:1px solid rgba(104,156,255,0.28);width:min(28rem,100%)">
                <input type="hidden" name="token" value="${token}">
                <h1 style="margin:0 0 1.5rem">Set New Password</h1>
                <label style="display:block;font-size:.85rem;color:#aab8d8;margin-bottom:.5rem">New Password</label>
                <input type="password" name="password" required placeholder="Enter new password" style="width:100%;padding:.75rem 1rem;border:1px solid rgba(169,199,197,0.3);border-radius:.5rem;background:rgba(5,14,20,0.8);color:#eff8f6;font-size:1rem;margin-bottom:1rem">
                <label style="display:block;font-size:.85rem;color:#aab8d8;margin-bottom:.5rem">Confirm Password</label>
                <input type="password" name="confirm" required placeholder="Confirm new password" style="width:100%;padding:.75rem 1rem;border:1px solid rgba(169,199,197,0.3);border-radius:.5rem;background:rgba(5,14,20,0.8);color:#eff8f6;font-size:1rem;margin-bottom:1rem">
                <button type="submit" style="width:100%;padding:.875rem;border:0;border-radius:.5rem;font-weight:700;color:#081021;background:linear-gradient(130deg,#3f7cff,#2cc8a7);cursor:pointer">Reset Password</button>
            </form>
        </body>
        </html>
    `, {
        headers: { 'Content-Type': 'text/html' }
    });
};
