import { execSync } from 'child_process';
import { describe, it, expect } from 'vitest';
import * as path from 'path';

// Note: This test requires the dev server to be running (`bun run dev`)
// or to be executed via `lhci autorun` which boots the server itself.
// This file serves as a wrapper to run Lighthouse assertions
// programmatically within the Vitest suite if desired.

describe('Lighthouse Performance', () => {
    it('passes performance budgets defined in lighthouserc.js', () => {
        try {
            // We run the Lighthouse CI CLI synchronously
            const rootDir = path.resolve(__dirname, '../../../');

            // `lhci autorun` reads from `lighthouserc.js`
            const result = execSync('npx lhci autorun', {
                cwd: rootDir,
                encoding: 'utf-8',
                stdio: 'pipe'
            });

            // If it exits with 0, it means all assertions passed.
            expect(result).toContain('Done running autorun');
        } catch (error: any) {
            // If Lighthouse assertions fail, lhci exits with non-zero
            console.error(error.stdout);
            console.error(error.stderr);
            throw new Error(`Lighthouse performance budgets failed:\n${error.message}`);
        }
    }, 120000); // 2 minute timeout for LH runs
});
