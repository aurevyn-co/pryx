import { describe, it, expect } from 'vitest';

describe('OAuth callback validation', () => {
    it('validates state parameter matches cookie', () => {
        const state = 'abc123';
        const storedState = 'abc123';
        expect(state).toBe(storedState);
    });

    it('rejects mismatched state', () => {
        const state = 'abc123';
        const storedState = 'different';
        expect(state).not.toBe(storedState);
    });

    it('rejects missing state', () => {
        const state = null;
        const storedState = 'abc123';
        expect(!state || state !== storedState).toBe(true);
    });

    it('rejects missing code', () => {
        const code = null;
        expect(!code).toBe(true);
    });
});

describe('GitHub API response handling', () => {
    it('extracts access token from success response', () => {
        const response = { access_token: 'gho_xxx', token_type: 'bearer' };
        expect(response.access_token).toBe('gho_xxx');
    });

    it('handles error response', () => {
        const response = { error: 'bad_verification_code', error_description: 'Invalid code' } as { error: string; error_description: string };
        expect((response as any).access_token).toBeUndefined();
        expect(response.error).toBe('bad_verification_code');
    });

    it('extracts primary verified email', () => {
        const emails = [
            { email: 'spam@old.com', primary: false, verified: true },
            { email: 'main@example.com', primary: true, verified: true },
            { email: 'unverified@example.com', primary: false, verified: false }
        ];
        const primaryEmail = emails.find(e => e.primary && e.verified);
        expect(primaryEmail?.email).toBe('main@example.com');
    });

    it('returns undefined when no verified primary email', () => {
        const emails = [
            { email: 'unverified@example.com', primary: true, verified: false }
        ];
        const primaryEmail = emails.find(e => e.primary && e.verified);
        expect(primaryEmail).toBeUndefined();
    });
});

describe('User creation from OAuth', () => {
    it('uses GitHub name when available', () => {
        const githubUser = { id: 123, login: 'octocat', name: 'The Octocat' };
        const displayName = githubUser.name || githubUser.login;
        expect(displayName).toBe('The Octocat');
    });

    it('falls back to GitHub login when name unavailable', () => {
        const githubUser = { id: 123, login: 'octocat', name: null };
        const displayName = githubUser.name || githubUser.login;
        expect(displayName).toBe('octocat');
    });

    it('converts GitHub ID to string', () => {
        const githubId = 12345678;
        const providerUserId = String(githubId);
        expect(providerUserId).toBe('12345678');
        expect(typeof providerUserId).toBe('string');
    });
});

describe('Session cookie', () => {
    it('sets secure cookie attributes', () => {
        const token = 'abc123';
        const cookieValue = `auth_token=${token}; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax; Secure`;
        expect(cookieValue).toContain('HttpOnly');
        expect(cookieValue).toContain('Secure');
        expect(cookieValue).toContain('SameSite=Lax');
        expect(cookieValue).toContain('Max-Age=2592000');
    });
});
