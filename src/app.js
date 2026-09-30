document.addEventListener('DOMContentLoaded', () => {
  const searchInput = document.getElementById('bundle-search');
  const bundleCards = document.querySelectorAll('.bundle-card');
  const filterPills = document.querySelectorAll('.filter-pill');
  const noResults = document.getElementById('no-results');

  let activeTier = 'all';

  function applyFilters() {
    const query = (searchInput?.value || '').toLowerCase().trim();
    let visibleCount = 0;

    bundleCards.forEach((card) => {
      const id = card.getAttribute('data-id')?.toLowerCase() || '';
      const title = card.getAttribute('data-title')?.toLowerCase() || '';
      const desc = card.getAttribute('data-desc')?.toLowerCase() || '';
      const tier = card.getAttribute('data-tier')?.toLowerCase() || 'official';

      const matchesQuery = !query || id.includes(query) || title.includes(query) || desc.includes(query);
      const matchesTier = activeTier === 'all' || tier === activeTier;

      if (matchesQuery && matchesTier) {
        card.style.display = 'flex';
        visibleCount++;
      } else {
        card.style.display = 'none';
      }
    });

    if (noResults) {
      noResults.style.display = visibleCount === 0 ? 'block' : 'none';
    }
  }

  if (searchInput) {
    searchInput.addEventListener('input', () => {
      applyFilters();
      const query = searchInput.value.trim();
      const newUrl = query
        ? `${window.location.pathname}?q=${encodeURIComponent(query)}`
        : window.location.pathname;
      window.history.replaceState({}, '', newUrl);
    });
  }

  filterPills.forEach((pill) => {
    pill.addEventListener('click', () => {
      filterPills.forEach((p) => p.classList.remove('active'));
      pill.classList.add('active');
      activeTier = pill.getAttribute('data-filter') || 'all';
      applyFilters();
    });
  });

  // URL Query Parameter Auto-Search (?q=... or ?search=...)
  const urlParams = new URLSearchParams(window.location.search);
  const initialQuery = urlParams.get('q') || urlParams.get('search');
  if (initialQuery && searchInput) {
    searchInput.value = initialQuery;
    activeTier = 'all';
    filterPills.forEach((p) => {
      p.classList.toggle('active', p.getAttribute('data-filter') === 'all');
    });
    applyFilters();

    // Smooth scroll and pulse highlight on first matching card
    setTimeout(() => {
      const match = Array.from(bundleCards).find((c) => c.style.display !== 'none');
      if (match) {
        match.scrollIntoView({ behavior: 'smooth', block: 'center' });
        match.classList.add('card-highlight');
        setTimeout(() => match.classList.remove('card-highlight'), 2200);
      }
    }, 150);
  }

  // Copy install command to clipboard handler
  document.querySelectorAll('.copy-btn').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const text = btn.getAttribute('data-copy');
      if (!text) return;
      try {
        await navigator.clipboard.writeText(text);
        const originalSvg = btn.innerHTML;
        btn.innerHTML = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`;
        setTimeout(() => {
          btn.innerHTML = originalSvg;
        }, 1800);
      } catch (err) {
        console.error('Failed to copy', err);
      }
    });
  });

  // Copy Badge markdown handler
  document.querySelectorAll('.badge-copy-btn').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const text = btn.getAttribute('data-badge');
      if (!text) return;
      try {
        await navigator.clipboard.writeText(text);
        const originalHtml = btn.innerHTML;
        btn.innerHTML = `<span style="color: #10b981;">✓ Copied!</span>`;
        setTimeout(() => {
          btn.innerHTML = originalHtml;
        }, 1800);
      } catch (err) {
        console.error('Failed to copy badge', err);
      }
    });
  });
});

