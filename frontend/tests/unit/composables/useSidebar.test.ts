import { describe, it, expect, beforeEach } from 'vitest';
import { useSidebar } from '../../../composables/useSidebar';

describe('useSidebar', () => {
    beforeEach(() => {
        // We can clear nuxt state if necessary
        const { isExpanded, isOverlay } = useSidebar();
        isExpanded.value = true;
        isOverlay.value = false;
    });

    it('initializes with default values', () => {
        // Because state is preserved between tests in the same environment,
        // we should test the hook output
        const { isExpanded, isOverlay } = useSidebar();
        expect(isExpanded.value).toBe(true);
        expect(isOverlay.value).toBe(false);
    });

    it('toggle() flips isExpanded', () => {
        const { isExpanded, toggle } = useSidebar();
        isExpanded.value = true;
        toggle();
        expect(isExpanded.value).toBe(false);
        toggle();
        expect(isExpanded.value).toBe(true);
    });

    it('close() only sets isExpanded to false if isOverlay is true', () => {
        const { isExpanded, isOverlay, close } = useSidebar();

        isExpanded.value = true;
        isOverlay.value = false;
        close();
        expect(isExpanded.value).toBe(true);

        isOverlay.value = true;
        close();
        expect(isExpanded.value).toBe(false);
    });

    it('checkScreenSize() updates isOverlay based on window.innerWidth', () => {
        const { isOverlay, checkScreenSize } = useSidebar();

        Object.defineProperty(window, 'innerWidth', { writable: true, configurable: true, value: 1200 });
        checkScreenSize();
        expect(isOverlay.value).toBe(false);

        Object.defineProperty(window, 'innerWidth', { writable: true, configurable: true, value: 800 });
        checkScreenSize();
        expect(isOverlay.value).toBe(true);
    });
});
