if (!api.isLoggedIn()) { window.location.href = '/login'; }

let pnlChartInstance = null;
let equityChartInstance = null;

// Sidebar mobile
document.getElementById('hamburgerBtn').addEventListener('click', () => {
    document.getElementById('sidebar').classList.add('open');
    document.getElementById('sidebarOverlay').classList.add('open');
});
document.getElementById('sidebarOverlay').addEventListener('click', () => {
    document.getElementById('sidebar').classList.remove('open');
    document.getElementById('sidebarOverlay').classList.remove('open');
});

// Alert bell
document.getElementById('alertBell').addEventListener('click', (e) => {
    e.stopPropagation();
    document.getElementById('alertDropdown').classList.toggle('open');
});
document.addEventListener('click', () => document.getElementById('alertDropdown').classList.remove('open'));
document.getElementById('alertDropdown').addEventListener('click', (e) => e.stopPropagation());

document.getElementById('checkAlertsBtn').addEventListener('click', async () => {
    try {
        await api.checkAlerts();
        await loadAlerts();
        Toast.success('اعلان‌ها بررسی شد');
    } catch (err) {
        Toast.error(err.message);
    }
});

async function loadAlerts() {
    try {
        const data = await api.getAlerts();
        const alerts = data.alerts || [];
        const countEl = document.getElementById('alertCount');
        const listEl = document.getElementById('alertList');

        countEl.textContent = alerts.length;
        countEl.dataset.count = alerts.length;

        if (alerts.length === 0) {
            listEl.innerHTML = '<div class="alert-empty">اعلانی وجود ندارد</div>';
            return;
        }
        listEl.innerHTML = alerts.map(a => `
            <div class="alert-item level-${a.level}" onclick="markAlertRead(${a.id})">
                <div>${a.message}</div>
                <div style="font-size:10px;color:var(--text-muted);margin-top:4px">${new Date(a.created_at).toLocaleString('fa-IR')}</div>
            </div>
        `).join('');
    } catch {}
}

async function markAlertRead(id) {
    try {
        await api.markAlertRead(id);
        await loadAlerts();
    } catch {}
}

// Calculator
let portfolioBalances = {};

document.getElementById('calcBtn').addEventListener('click', async () => {
    const portfolioID = parseInt(document.getElementById('calcPortfolio').value);
    const risk = parseFloat(document.getElementById('calcRisk').value);
    const leverage = parseFloat(document.getElementById('calcLeverage').value) || 1;
    const entry = parseFloat(document.getElementById('calcEntry').value);
    const sl = parseFloat(document.getElementById('calcSL').value);
    const tp = parseFloat(document.getElementById('calcTP').value) || null;

    Validate.clearAll(document.querySelector('.calculator'));
    let hasError = false;
    if (!portfolioID) { Toast.warning('پورتفولیو را انتخاب کنید'); hasError = true; }
    if (!entry) { Validate.showFieldError(document.getElementById('calcEntry'), 'قیمت ورود الزامی است'); hasError = true; }
    if (!sl) { Validate.showFieldError(document.getElementById('calcSL'), 'حد ضرر الزامی است'); hasError = true; }
    if (hasError) return;

    try {
        const data = await api.calculate(portfolioID, entry, sl, tp, risk);
        const r = data.result;
        const posVal = r.position_value;
        const margin = posVal / leverage;
        const qty = r.quantity;

        document.getElementById('resRiskAmt').textContent = '$' + r.risk_amount.toFixed(2);
        document.getElementById('resQty').textContent = qty.toFixed(5);
        document.getElementById('resMargin').textContent = '$' + margin.toLocaleString('en', { maximumFractionDigits: 2 });
        document.getElementById('resPosVal').textContent = '$' + posVal.toLocaleString('en', { maximumFractionDigits: 2 });
        document.getElementById('resRR').textContent = r.risk_reward > 0 ? '1 : ' + r.risk_reward.toFixed(2) : '---';
        document.getElementById('calcResult').style.display = 'block';

        const balance = portfolioBalances[portfolioID] || 0;
        const warningEl = document.getElementById('calcWarning');
        if (margin > balance && balance > 0) {
            warningEl.style.display = 'flex';
        } else {
            warningEl.style.display = 'none';
        }

        Toast.success('محاسبه انجام شد');
    } catch (err) {
        Toast.error(err.message);
    }
});

// Logout
document.getElementById('logoutBtn').addEventListener('click', () => api.logout());

// Keyboard shortcuts
Shortcuts.register('b', false, false, () => document.getElementById('calcEntry').focus());
Shortcuts.init();

// Load all in parallel
(async () => {
    try {
        const [dashResult, pnlResult, tradesResult, portfoliosResult] = await Promise.all([
            api.getDashboard().catch(() => null),
            api.getPnLHistory().catch(() => null),
            api.getTrades(5, 0).catch(() => null),
            api.getPortfolios().catch(() => null),
        ]);

        if (dashResult) {
            const d = dashResult.dashboard;
            const pnlPos = d.total_pnl >= 0;
            document.getElementById('statsGrid').innerHTML = `
                <div class="stat-card">
                    <div class="stat-label">موجودی کل</div>
                    <div class="stat-value" data-count="${d.balance}" data-prefix="$">---</div>
                    <div class="stat-sub"></div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">کل PnL</div>
                    <div class="stat-value ${pnlPos ? 'up' : 'down'}" data-count="${d.total_pnl}" data-prefix="${pnlPos ? '+$' : '$'}" data-decimals="0">---</div>
                    <div class="stat-sub">${d.closed_trades} ترید بسته شده</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">نرخ برد</div>
                    <div class="stat-value gold" data-count="${d.win_rate}" data-suffix="%" data-decimals="1">---</div>
                    <div class="stat-sub">بهترین: +$${d.best_trade.toFixed(0)} | بدترین: $${d.worst_trade.toFixed(0)}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">تریدهای باز</div>
                    <div class="stat-value" data-count="${d.open_trades}">---</div>
                    <div class="stat-sub">میانگین R/R: ${d.avg_risk_reward.toFixed(2)}</div>
                </div>
            `;
            Counter.animateAll(document.getElementById('statsGrid'));

            if (pnlResult) {
                const points = pnlResult.pnl_history;
                if (points && points.length > 0) {
                    const labels = points.map(p => new Date(p.date).toLocaleDateString('fa-IR'));
                    const values = points.map(p => p.pnl);

                    if (pnlChartInstance) pnlChartInstance.destroy();
                    const ctx = document.getElementById('pnlChart').getContext('2d');
                    pnlChartInstance = new Chart(ctx, {
                        type: 'line',
                        data: {
                            labels,
                            datasets: [{
                                data: values,
                                borderColor: '#C9A84C',
                                backgroundColor: 'rgba(201, 168, 76, 0.1)',
                                borderWidth: 2, pointRadius: 3, pointBackgroundColor: '#C9A84C',
                                fill: true, tension: 0.4
                            }]
                        },
                        options: {
                            responsive: true, maintainAspectRatio: false,
                            plugins: { legend: { display: false } },
                            scales: {
                                x: { ticks: { color: '#555', font: { family: 'Vazirmatn', size: 10 } }, grid: { color: '#ffffff08' } },
                                y: { ticks: { color: '#555', font: { family: 'Vazirmatn', size: 10 } }, grid: { color: '#ffffff08' } }
                            }
                        }
                    });

                    const equityPoints = points.map((p, i) => ({
                        date: p.date,
                        value: d.balance - (points[points.length - 1].pnl) + p.pnl
                    }));
                    const eLabels = equityPoints.map(p => new Date(p.date).toLocaleDateString('fa-IR'));
                    const eValues = equityPoints.map(p => p.value);

                    if (equityChartInstance) equityChartInstance.destroy();
                    const ectx = document.getElementById('equityChart').getContext('2d');
                    equityChartInstance = new Chart(ectx, {
                        type: 'line',
                        data: {
                            labels: eLabels,
                            datasets: [{
                                data: eValues,
                                borderColor: '#3B82F6',
                                backgroundColor: 'rgba(59, 130, 246, 0.1)',
                                borderWidth: 2, pointRadius: 3, pointBackgroundColor: '#3B82F6',
                                fill: true, tension: 0.4
                            }]
                        },
                        options: {
                            responsive: true, maintainAspectRatio: false,
                            plugins: { legend: { display: false } },
                            scales: {
                                x: { ticks: { color: '#555', font: { family: 'Vazirmatn', size: 10 } }, grid: { color: '#ffffff08' } },
                                y: { ticks: { color: '#555', font: { family: 'Vazirmatn', size: 10 }, callback: v => '$' + v.toLocaleString() }, grid: { color: '#ffffff08' } }
                            }
                        }
                    });
                } else {
                    document.getElementById('chartEmpty').style.display = 'block';
                    document.getElementById('equityEmpty').style.display = 'block';
                }
            }
        }

        if (tradesResult) {
            const trades = tradesResult.trades;
            const list = document.getElementById('tradesList');
            if (!trades || trades.length === 0) { list.innerHTML = EmptyState.trades(); }
            else {
                list.innerHTML = trades.map(t => {
                    const pnlClass = t.pnl >= 0 ? 'up' : 'down';
                    const pnlText = t.pnl != null ? (t.pnl >= 0 ? '+' : '') + '$' + t.pnl.toFixed(0) : '---';
                    return `<div class="trade-row">
                        <div>
                            <div style="display:flex;align-items:center;flex-direction:row-reverse">
                                <span class="trade-sym">${t.symbol}</span>
                                <span class="trade-badge ${t.side === 'LONG' ? 'long' : 'short'}">${t.side}</span>
                            </div>
                            <div class="trade-info">ورود: $${t.entry_price} · SL: $${t.stop_loss}</div>
                        </div>
                        <div>
                            <div class="trade-pnl ${pnlClass}">${pnlText}</div>
                            <div style="text-align:left;margin-top:4px"><span class="status-badge">${t.status}</span></div>
                        </div>
                    </div>`;
                }).join('');
            }
        }

        if (portfoliosResult) {
            const portfolios = portfoliosResult.portfolios || [];
            const select = document.getElementById('calcPortfolio');
            portfolioBalances = {};
            select.innerHTML = portfolios.length > 0
                ? portfolios.map(p => {
                    portfolioBalances[p.ID] = p.balance;
                    return `<option value="${p.ID}">${p.name} - $${p.balance.toLocaleString('en')}</option>`;
                }).join('')
                : '<option value="">ابتدا پورتفولیو بسازید</option>';
        }

        loadAlerts();
    } catch {}
})();
