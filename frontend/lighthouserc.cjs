module.exports = {
    ci: {
        collect: {
            url: ['http://localhost:3000/login'],
            startServerCommand: 'bun run build && PORT=3000 bun run preview',
            startServerReadyPattern: 'Listening on http://',
            numberOfRuns: 1,
        },
        assert: {
            assertions: {
                'categories:performance': ['warn', { minScore: 0.7 }],
                'first-contentful-paint': ['warn', { maxNumericValue: 2500 }],
                'largest-contentful-paint': ['warn', { maxNumericValue: 3000 }],
                'interactive': ['warn', { maxNumericValue: 4000 }],
            },
        },
        upload: {
            target: 'temporary-public-storage',
        },
    },
};
