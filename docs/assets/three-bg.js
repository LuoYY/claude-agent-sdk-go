/**
 * Subtle Three.js particle network background
 * White theme with light gray interconnected particles
 */
(function () {
    const canvas = document.getElementById('bg-canvas');
    if (!canvas || typeof THREE === 'undefined') return;

    // Scene setup
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x0d1117);

    const camera = new THREE.PerspectiveCamera(
        60,
        window.innerWidth / window.innerHeight,
        0.1,
        1000
    );
    camera.position.z = 50;

    const renderer = new THREE.WebGLRenderer({
        canvas: canvas,
        antialias: true,
        alpha: false,
    });
    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

    // Mouse tracking
    const mouse = { x: 0, y: 0, target: { x: 0, y: 0 } };

    // Particle system
    const PARTICLE_COUNT = 120;
    const SPREAD = 80;
    const CONNECTION_DISTANCE = 12;
    const PARTICLE_SIZE = 0.15;

    // Create particles
    const particleGeometry = new THREE.BufferGeometry();
    const positions = new Float32Array(PARTICLE_COUNT * 3);
    const velocities = [];
    const opacities = new Float32Array(PARTICLE_COUNT);

    for (let i = 0; i < PARTICLE_COUNT; i++) {
        positions[i * 3] = (Math.random() - 0.5) * SPREAD;
        positions[i * 3 + 1] = (Math.random() - 0.5) * SPREAD;
        positions[i * 3 + 2] = (Math.random() - 0.5) * 30;

        velocities.push({
            x: (Math.random() - 0.5) * 0.015,
            y: (Math.random() - 0.5) * 0.015,
            z: (Math.random() - 0.5) * 0.005,
        });

        opacities[i] = 0.15 + Math.random() * 0.25;
    }

    particleGeometry.setAttribute(
        'position',
        new THREE.BufferAttribute(positions, 3)
    );

    // Custom shader for particles with varying opacity
    const particleMaterial = new THREE.ShaderMaterial({
        uniforms: {
            uTime: { value: 0 },
            uPixelRatio: { value: renderer.getPixelRatio() },
        },
        vertexShader: `
            uniform float uTime;
            uniform float uPixelRatio;
            attribute float opacity;
            varying float vOpacity;

            void main() {
                vOpacity = opacity;
                vec4 mvPosition = modelViewMatrix * vec4(position, 1.0);
                gl_PointSize = ${PARTICLE_SIZE.toFixed(1)} * uPixelRatio * (80.0 / -mvPosition.z);
                gl_Position = projectionMatrix * mvPosition;
            }
        `,
        fragmentShader: `
            varying float vOpacity;

            void main() {
                float dist = length(gl_PointCoord - vec2(0.5));
                if (dist > 0.5) discard;

                float alpha = smoothstep(0.5, 0.1, dist) * vOpacity;
                gl_FragColor = vec4(0.35, 0.40, 0.50, alpha);
            }
        `,
        transparent: true,
        depthWrite: false,
    });

    particleGeometry.setAttribute(
        'opacity',
        new THREE.BufferAttribute(opacities, 1)
    );

    const particles = new THREE.Points(particleGeometry, particleMaterial);
    scene.add(particles);

    // Connection lines
    const MAX_LINES = PARTICLE_COUNT * 6;
    const lineGeometry = new THREE.BufferGeometry();
    const linePositions = new Float32Array(MAX_LINES * 6);
    const lineColors = new Float32Array(MAX_LINES * 8);

    lineGeometry.setAttribute(
        'position',
        new THREE.BufferAttribute(linePositions, 3)
    );
    lineGeometry.setAttribute(
        'color',
        new THREE.BufferAttribute(lineColors, 4)
    );

    const lineMaterial = new THREE.LineBasicMaterial({
        vertexColors: true,
        transparent: true,
        depthWrite: false,
        blending: THREE.NormalBlending,
    });

    // Override to use vec4 colors
    lineMaterial.onBeforeCompile = (shader) => {
        shader.vertexShader = shader.vertexShader.replace(
            'attribute vec3 color;',
            'attribute vec4 color;'
        );
        shader.vertexShader = shader.vertexShader.replace(
            'vColor = color;',
            'vColor = color;'
        );
    };

    const lines = new THREE.LineSegments(lineGeometry, lineMaterial);
    scene.add(lines);

    // Floating geometric shapes (very subtle)
    const shapes = [];
    const shapeMaterial = new THREE.MeshBasicMaterial({
        color: 0x30363d,
        transparent: true,
        opacity: 0.3,
        wireframe: true,
    });

    for (let i = 0; i < 5; i++) {
        let geometry;
        const type = Math.floor(Math.random() * 3);
        if (type === 0) {
            geometry = new THREE.IcosahedronGeometry(2 + Math.random() * 3, 0);
        } else if (type === 1) {
            geometry = new THREE.OctahedronGeometry(2 + Math.random() * 2, 0);
        } else {
            geometry = new THREE.TetrahedronGeometry(2 + Math.random() * 2, 0);
        }

        const mesh = new THREE.Mesh(geometry, shapeMaterial.clone());
        mesh.position.set(
            (Math.random() - 0.5) * 60,
            (Math.random() - 0.5) * 40,
            -10 - Math.random() * 20
        );
        mesh.rotation.set(
            Math.random() * Math.PI,
            Math.random() * Math.PI,
            Math.random() * Math.PI
        );
        mesh.userData = {
            rotSpeed: {
                x: (Math.random() - 0.5) * 0.003,
                y: (Math.random() - 0.5) * 0.003,
                z: (Math.random() - 0.5) * 0.002,
            },
            floatSpeed: 0.0003 + Math.random() * 0.0005,
            floatOffset: Math.random() * Math.PI * 2,
        };
        scene.add(mesh);
        shapes.push(mesh);
    }

    // Animation
    let time = 0;
    let lineIndex = 0;

    function animate() {
        requestAnimationFrame(animate);
        time += 0.016;

        // Smooth mouse tracking
        mouse.x += (mouse.target.x - mouse.x) * 0.05;
        mouse.y += (mouse.target.y - mouse.y) * 0.05;

        // Update particles
        const pos = particleGeometry.attributes.position.array;

        for (let i = 0; i < PARTICLE_COUNT; i++) {
            const ix = i * 3;
            const iy = i * 3 + 1;
            const iz = i * 3 + 2;

            pos[ix] += velocities[i].x;
            pos[iy] += velocities[i].y;
            pos[iz] += velocities[i].z;

            // Gentle mouse influence
            pos[ix] += mouse.x * 0.0003;
            pos[iy] += mouse.y * 0.0003;

            // Boundaries (wrap around)
            const halfSpread = SPREAD / 2;
            if (pos[ix] > halfSpread) pos[ix] = -halfSpread;
            if (pos[ix] < -halfSpread) pos[ix] = halfSpread;
            if (pos[iy] > halfSpread) pos[iy] = -halfSpread;
            if (pos[iy] < -halfSpread) pos[iy] = halfSpread;
            if (pos[iz] > 15) pos[iz] = -15;
            if (pos[iz] < -15) pos[iz] = 15;
        }

        particleGeometry.attributes.position.needsUpdate = true;

        // Update connections
        lineIndex = 0;
        const lp = lineGeometry.attributes.position.array;
        const lc = lineGeometry.attributes.color.array;

        for (let i = 0; i < PARTICLE_COUNT && lineIndex < MAX_LINES; i++) {
            for (let j = i + 1; j < PARTICLE_COUNT && lineIndex < MAX_LINES; j++) {
                const dx = pos[i * 3] - pos[j * 3];
                const dy = pos[i * 3 + 1] - pos[j * 3 + 1];
                const dz = pos[i * 3 + 2] - pos[j * 3 + 2];
                const dist = Math.sqrt(dx * dx + dy * dy + dz * dz);

                if (dist < CONNECTION_DISTANCE) {
                    const alpha = (1 - dist / CONNECTION_DISTANCE) * 0.08;
                    const li = lineIndex * 6;
                    const ci = lineIndex * 8;

                    lp[li] = pos[i * 3];
                    lp[li + 1] = pos[i * 3 + 1];
                    lp[li + 2] = pos[i * 3 + 2];
                    lp[li + 3] = pos[j * 3];
                    lp[li + 4] = pos[j * 3 + 1];
                    lp[li + 5] = pos[j * 3 + 2];

                    lc[ci] = 0.35;     lc[ci + 1] = 0.40;  lc[ci + 2] = 0.50;   lc[ci + 3] = alpha;
                    lc[ci + 4] = 0.35; lc[ci + 5] = 0.40;  lc[ci + 6] = 0.50;   lc[ci + 7] = alpha;

                    lineIndex++;
                }
            }
        }

        // Clear remaining lines
        for (let i = lineIndex; i < MAX_LINES; i++) {
            const li = i * 6;
            lp[li] = lp[li + 1] = lp[li + 2] = 0;
            lp[li + 3] = lp[li + 4] = lp[li + 5] = 0;
        }

        lineGeometry.attributes.position.needsUpdate = true;
        lineGeometry.attributes.color.needsUpdate = true;
        lineGeometry.setDrawRange(0, lineIndex * 2);

        // Update geometric shapes
        shapes.forEach((shape) => {
            shape.rotation.x += shape.userData.rotSpeed.x;
            shape.rotation.y += shape.userData.rotSpeed.y;
            shape.rotation.z += shape.userData.rotSpeed.z;
            shape.position.y +=
                Math.sin(time * shape.userData.floatSpeed * 60 + shape.userData.floatOffset) * 0.01;
        });

        // Subtle camera movement
        camera.position.x = mouse.x * 0.5;
        camera.position.y = mouse.y * 0.3;
        camera.lookAt(0, 0, 0);

        particleMaterial.uniforms.uTime.value = time;
        renderer.render(scene, camera);
    }

    animate();

    // Event listeners
    document.addEventListener('mousemove', (e) => {
        mouse.target.x = (e.clientX / window.innerWidth - 0.5) * 2;
        mouse.target.y = -(e.clientY / window.innerHeight - 0.5) * 2;
    });

    window.addEventListener('resize', () => {
        camera.aspect = window.innerWidth / window.innerHeight;
        camera.updateProjectionMatrix();
        renderer.setSize(window.innerWidth, window.innerHeight);
        renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
        particleMaterial.uniforms.uPixelRatio.value = renderer.getPixelRatio();
    });

    // Reduce animation when tab not visible
    document.addEventListener('visibilitychange', () => {
        if (document.hidden) {
            renderer.setAnimationLoop(null);
        } else {
            renderer.setAnimationLoop(null);
            animate();
        }
    });
})();
