import { describe, expect, it } from 'vitest';
import { apiApp } from './[...path]';

describe('API routing', () => {
    it('routes /api/admin/telemetry/config through admin router', async () => {
        const response = await apiApp.request('http://localhost/api/admin/telemetry/config');

        expect(response.status).toBe(200);
        const body = await response.json();
        expect(body).toHaveProperty('enabled');
    });
});
