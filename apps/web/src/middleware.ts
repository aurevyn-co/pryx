import { defineMiddleware } from 'astro:middleware';

export const onRequest = defineMiddleware(async (context, next) => {
    // Ensure locals.runtime.env is available and contains bindings for Hono
    if (context.locals) {
        const platformEnv = (context as any).platform?.env || {};

        // Set up locals.runtime.env with all bindings for Hono to access
        context.locals.runtime = {
            env: platformEnv,
            cf: (context as any).cf,
            caches: (context as any).caches,
            ctx: {
                waitUntil: (promise: Promise<any>) => context.waitUntil(promise),
                passThroughOnException: () => {},
            },
        };
    }
    return next();
});
