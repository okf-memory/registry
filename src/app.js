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
    return visibleCount;
  }

  let searchDebounceTimer = null;
  if (searchInput) {
    searchInput.addEventListener('input', () => {
      const visibleCount = applyFilters();
      const query = searchInput.value.trim();
      const newUrl = query
        ? `${window.location.pathname}?q=${encodeURIComponent(query)}`
        : window.location.pathname;
      window.history.replaceState({}, '', newUrl);

      clearTimeout(searchDebounceTimer);
      if (query.length >= 2) {
        searchDebounceTimer = setTimeout(() => {
          if (typeof gtag === 'function') {
            gtag('event', 'registry_search', {
              'search_term': query.toLowerCase(),
              'result_count': visibleCount
            });
          }
        }, 500);
      }
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
    const visibleCount = applyFilters();

    if (typeof gtag === 'function' && initialQuery.length >= 2) {
      gtag('event', 'registry_search', {
        'search_term': initialQuery.toLowerCase(),
        'result_count': visibleCount
      });
    }

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

        // GA4 Tracking
        const card = btn.closest('.bundle-card');
        const bundleId = card?.getAttribute('data-id') || 'canonical_badge';
        const bundleTier = card?.getAttribute('data-tier') || 'official';
        const isBadge = text.includes('badge') || text.includes('.svg');

        if (typeof gtag === 'function') {
          if (isBadge) {
            gtag('event', 'badge_markdown_copy', {
              'badge_id': bundleId
            });
          } else {
            gtag('event', 'bundle_install_copy', {
              'bundle_id': bundleId,
              'bundle_tier': bundleTier
            });
          }
        }
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

        // GA4 Tracking
        const card = btn.closest('.bundle-card');
        const bundleId = card?.getAttribute('data-id') || 'badge';
        if (typeof gtag === 'function') {
          gtag('event', 'badge_markdown_copy', {
            'badge_id': bundleId
          });
        }
      } catch (err) {
        console.error('Failed to copy badge', err);
      }
    });
  });

  // GitHub Outbound Tracking
  document.querySelectorAll('a[href*="github.com"]').forEach(link => {
    link.addEventListener('click', () => {
      let repo = 'registry';
      const href = link.href || '';
      if (href.includes('okf-agent-memory')) repo = 'okf-agent-memory';
      else if (href.includes('topics/')) repo = 'topics/okf-memory-bundle';

      let placement = 'body';
      if (link.closest('header') || link.closest('.nav-wrap')) placement = 'nav';
      else if (link.closest('footer')) placement = 'footer';
      else if (link.closest('.submit-action') || link.closest('#publish-bundle')) placement = 'publish';

      if (typeof gtag === 'function') {
        gtag('event', 'github_outbound_click', {
          'target_repo': repo,
          'placement': placement
        });
      }
    });
  });

  // ------------------------------------------------------------------------
  // DSGVO / GDPR Cookie Consent Manager (Google Consent Mode v2)
  // ------------------------------------------------------------------------
  const cookieBanner = document.getElementById('cookie-banner');
  const cookieBtnAccept = document.getElementById('cookie-btn-accept');
  const cookieBtnReject = document.getElementById('cookie-btn-reject');
  const cookieSettingsBtn = document.getElementById('cookie-settings-btn');

  function getConsentCookie() {
    const match = document.cookie.match(/(^|;)\s*okf_cookie_consent\s*=\s*([^;]+)/);
    if (match) return decodeURIComponent(match[2]);
    try { return localStorage.getItem('okf_cookie_consent'); } catch (e) {}
    return null;
  }

  function setConsentCookie(value) {
    const maxAge = 365 * 24 * 60 * 60; // 1 year
    let cookieStr = `okf_cookie_consent=${encodeURIComponent(value)}; path=/; max-age=${maxAge}; SameSite=Lax`;
    const hostname = window.location.hostname;
    if (hostname === 'okf-memory.dev' || hostname.endsWith('.okf-memory.dev')) {
      cookieStr += '; domain=.okf-memory.dev';
    }
    document.cookie = cookieStr;
    try { localStorage.setItem('okf_cookie_consent', value); } catch (e) {}
  }

  function openCookieBanner() {
    if (!cookieBanner) return;
    cookieBanner.style.display = 'block';
    requestAnimationFrame(() => {
      cookieBanner.classList.add('show');
    });
  }

  function closeCookieBanner() {
    if (!cookieBanner) return;
    cookieBanner.classList.remove('show');
    setTimeout(() => {
      if (!cookieBanner.classList.contains('show')) {
        cookieBanner.style.display = 'none';
      }
    }, 350);
  }

  const existingConsent = getConsentCookie();
  if (!existingConsent) {
    setTimeout(openCookieBanner, 600);
  }

  if (cookieBtnAccept) {
    cookieBtnAccept.addEventListener('click', () => {
      setConsentCookie('granted');
      if (typeof gtag === 'function') {
        gtag('consent', 'update', {
          'analytics_storage': 'granted'
        });
        gtag('event', 'page_view', {
          page_title: document.title,
          page_location: window.location.href
        });
      }
      closeCookieBanner();
    });
  }

  if (cookieBtnReject) {
    cookieBtnReject.addEventListener('click', () => {
      setConsentCookie('denied');
      if (typeof gtag === 'function') {
        gtag('consent', 'update', {
          'analytics_storage': 'denied'
        });
      }
      closeCookieBanner();
    });
  }

  if (cookieSettingsBtn) {
    cookieSettingsBtn.addEventListener('click', (e) => {
      e.preventDefault();
      openCookieBanner();
    });
  }
});


