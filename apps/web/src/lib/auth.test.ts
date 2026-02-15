import { describe, it, expect, beforeEach } from 'vitest';
import {
    isValidEmail,
    isValidPassword,
    generateId,
    generateToken,
    hashToken
} from '../lib/auth';

describe('Auth validation', () => {
    describe('isValidEmail', () => {
        it('accepts valid emails', () => {
            expect(isValidEmail('test@example.com')).toBe(true);
            expect(isValidEmail('user.name+tag@domain.co')).toBe(true);
            expect(isValidEmail('a@b.c')).toBe(true);
        });

        it('rejects invalid emails', () => {
            expect(isValidEmail('')).toBe(false);
            expect(isValidEmail('test')).toBe(false);
            expect(isValidEmail('test@')).toBe(false);
            expect(isValidEmail('@domain.com')).toBe(false);
            expect(isValidEmail('test@domain')).toBe(false);
            expect(isValidEmail('test domain.com')).toBe(false);
        });
    });

    describe('isValidPassword', () => {
        it('accepts strong passwords', () => {
            expect(isValidPassword('Password1')).toEqual({ valid: true });
            expect(isValidPassword('Abcdefg1')).toEqual({ valid: true });
            expect(isValidPassword('MyP@ssw0rd')).toEqual({ valid: true });
        });

        it('rejects short passwords', () => {
            const result = isValidPassword('Pass1');
            expect(result.valid).toBe(false);
            expect(result.error).toContain('8 characters');
        });

        it('rejects passwords without uppercase', () => {
            const result = isValidPassword('password1');
            expect(result.valid).toBe(false);
            expect(result.error).toContain('uppercase');
        });

        it('rejects passwords without lowercase', () => {
            const result = isValidPassword('PASSWORD1');
            expect(result.valid).toBe(false);
            expect(result.error).toContain('lowercase');
        });

        it('rejects passwords without number', () => {
            const result = isValidPassword('Passwords');
            expect(result.valid).toBe(false);
            expect(result.error).toContain('number');
        });
    });
});

describe('Token generation', () => {
    it('generates unique IDs', () => {
        const id1 = generateId();
        const id2 = generateId();
        expect(id1).not.toBe(id2);
        expect(id1).toMatch(/^[0-9a-f-]{36}$/);
    });

    it('generates unique tokens', () => {
        const token1 = generateToken();
        const token2 = generateToken();
        expect(token1).not.toBe(token2);
        expect(token1).toHaveLength(64);
    });

    it('hashes tokens consistently', () => {
        const token = 'abc123';
        const hash1 = hashToken(token);
        const hash2 = hashToken(token);
        expect(hash1).toBe(hash2);
        expect(hash1).toHaveLength(16);
    });

    it('produces different hashes for different tokens', () => {
        const hash1 = hashToken('token1');
        const hash2 = hashToken('token2');
        expect(hash1).not.toBe(hash2);
    });
});
