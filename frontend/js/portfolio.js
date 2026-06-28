if (!api.isLoggedIn()) { window.location.href = '/login'; }

// Sidebar mobile
document.getElementById('hamburgerBtn').addEventListener('click', () => {
    document.getElementById('sidebar').classList.add('open');
    document.getElementById('sidebarOverlay').classList.add('open');
});
document.getElementById('sidebarOverlay').addEventListener('click', () => {
    document.getElementById('sidebar').classList.remove('open');
    document.getElementById('sidebarOverlay').classList.remove('open');
});

const loadPortfolios = async () => {
    const grid = document.getElementById('portfolioGrid');
    grid.innerHTML = Skeleton.rows(3);
    try {
        const data = await api.getPortfolios();
        const portfolios = data.portfolios || [];

        let totalCapital = 0, totalBalance = 0;
        portfolios.forEach(p => { totalCapital += p.capital; totalBalance += p.balance; });
        const totalPnl = totalBalance - totalCapital;
        const totalReturn = totalCapital > 0 ? ((totalPnl / totalCapital) * 100) : 0;

        document.getElementById('summaryGrid').innerHTML = `
            <div class="sum-card"><div class="sum-label">کل سرمایه</div><div class="sum-val gold">$${totalCapital.toLocaleString('en')}</div></div>
            <div class="sum-card"><div class="sum-label">کل موجودی</div><div class="sum-val">$${totalBalance.toLocaleString('en')}</div></div>
            <div class="sum-card"><div class="sum-label">کل سود/ضرر</div><div class="sum-val ${totalPnl >= 0 ? 'up' : 'down'}">${totalPnl >= 0 ? '+' : ''}$${totalPnl.toFixed(2)}</div></div>
            <div class="sum-card"><div class="sum-label">بازدهی کل</div><div class="sum-val ${totalReturn >= 0 ? 'up' : 'down'}">${totalReturn >= 0 ? '+' : ''}${totalReturn.toFixed(1)}%</div></div>
        `;

        if (portfolios.length === 0) {
            grid.innerHTML = EmptyState.portfolios();
            return;
        }

        grid.innerHTML = portfolios.map(p => {
            const pnl = p.balance - p.capital;
            const pnlPos = pnl >= 0;
            const ret = p.capital > 0 ? ((pnl / p.capital) * 100) : 0;
            const progress = Math.min((p.balance / p.capital) * 100, 200);
            const progressWidth = Math.min(progress, 100);
            const circumference = 2 * Math.PI * 30;
            const dashArray = (Math.min(Math.abs(ret), 100) / 100) * circumference;
            const donutColor = pnlPos ? '#C9A84C' : '#EF4444';

            return `<div class="pt-card ${p.is_default ? 'default' : ''}">
                <div class="pt-header"><div>
                    <div class="pt-name">${p.name}</div>
                    <div class="pt-trade-count">پورتفولیو</div>
                </div>
                <div style="display:flex;flex-direction:column;align-items:flex-end;gap:8px">
                    ${p.is_default ? '<span class="pt-default-badge"><i class="ti ti-star" style="font-size:10px"></i> پیشفرض</span>' : '<div style="height:26px"></div>'}
                    <div class="pt-donut">
                        <svg viewBox="0 0 80 80">
                            <circle cx="40" cy="40" r="30" fill="none" stroke="#1E1E1E" stroke-width="8"/>
                            <circle cx="40" cy="40" r="30" fill="none" stroke="${donutColor}" stroke-width="8" stroke-dasharray="${dashArray.toFixed(0)} ${circumference.toFixed(0)}" stroke-linecap="round"/>
                        </svg>
                        <div class="pt-donut-val" style="color:${donutColor}">${ret.toFixed(0)}%</div>
                    </div>
                </div></div>
                <div class="pt-stats">
                    <div><div class="pt-stat-label">سرمایه اولیه</div><div class="pt-stat-val gold">$${p.capital.toLocaleString('en')}</div></div>
                    <div><div class="pt-stat-label">موجودی فعلی</div><div class="pt-stat-val">$${p.balance.toLocaleString('en')}</div></div>
                    <div><div class="pt-stat-label">سود/ضرر</div><div class="pt-stat-val ${pnlPos ? 'up' : 'down'}">${pnlPos ? '+' : ''}$${pnl.toFixed(2)}</div></div>
                    <div><div class="pt-stat-label">بازدهی</div><div class="pt-stat-val ${pnlPos ? 'up' : 'down'}">${pnlPos ? '+' : ''}${ret.toFixed(1)}%</div></div>
                </div>
                <div class="pt-progress">
                    <div class="pt-progress-top"><span>رشد سرمایه</span><span style="color:${pnlPos ? '#C9A84C' : '#EF4444'}">${progress.toFixed(1)}%</span></div>
                    <div class="pt-bar"><div class="pt-fill" style="width:${progressWidth}%;background:${pnlPos ? 'linear-gradient(90deg,#8B6914,#C9A84C)' : 'linear-gradient(90deg,#7f1d1d,#EF4444)'}"></div></div>
                </div>
                <div class="pt-divider"></div>
                <div class="pt-actions">
                    ${!p.is_default ? `<button class="btn btn-ghost btn-sm" onclick="setDefault(${p.ID})"><i class="ti ti-star"></i> پیشفرض</button>` : ''}
                    <button class="btn btn-ghost btn-sm" onclick="window.location.href='/trades'"><i class="ti ti-chart-bar"></i> تریدها</button>
                    ${!p.is_default ? `<button class="btn btn-danger btn-sm" onclick="deletePortfolio(${p.ID})"><i class="ti ti-trash"></i></button>` : ''}
                </div>
            </div>`;
        }).join('');

    } catch (err) {
        grid.innerHTML = EmptyState.portfolios();
    }
};

const setDefault = async (id) => {
    try {
        await api.setDefaultPortfolio(id);
        Toast.success('پورتفولیو پیشفرض شد');
        loadPortfolios();
    } catch (err) { Toast.error(err.message); }
};

const deletePortfolio = async (id) => {
    const ok = await Confirm.show('حذف پورتفولیو', 'آیا از حذف این پورتفولیو اطمینان دارید؟');
    if (!ok) return;
    try {
        await api.deletePortfolio(id);
        Toast.success('پورتفولیو حذف شد');
        loadPortfolios();
    } catch (err) { Toast.error(err.message); }
};

document.getElementById('newPtBtn').addEventListener('click', () => {
    document.getElementById('modalOverlay').classList.add('open');
});
document.getElementById('modalClose').addEventListener('click', () => {
    document.getElementById('modalOverlay').classList.remove('open');
});

document.getElementById('submitPortfolio').addEventListener('click', async () => {
    Validate.clearAll(document.getElementById('modalOverlay'));
    const name = document.getElementById('ptName').value.trim();
    const capital = parseFloat(document.getElementById('ptCapital').value);

    let hasError = false;
    if (!name) { Validate.showFieldError(document.getElementById('ptName'), 'نام الزامی است'); hasError = true; }
    if (!capital || capital <= 0) { Validate.showFieldError(document.getElementById('ptCapital'), 'سرمایه باید بیشتر از صفر باشد'); hasError = true; }
    if (hasError) return;

    try {
        await api.createPortfolio(name, capital);
        document.getElementById('modalOverlay').classList.remove('open');
        document.getElementById('ptName').value = '';
        document.getElementById('ptCapital').value = '';
        Toast.success('پورتفولیو ایجاد شد');
        loadPortfolios();
    } catch (err) { Toast.error(err.message); }
});

document.getElementById('logoutBtn').addEventListener('click', () => api.logout());

Shortcuts.init();
loadPortfolios();
