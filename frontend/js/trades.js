if (!api.isLoggedIn()) { window.location.href = '/login'; }

let currentPage = 0;
let selectedTradeID = null;
let allTrades = [];
const PAGE_SIZE = 10;

// Checklist state
let checklistDefaults = [];
let currentChecklist = [];

// Sidebar mobile
document.getElementById('hamburgerBtn').addEventListener('click', () => {
    document.getElementById('sidebar').classList.add('open');
    document.getElementById('sidebarOverlay').classList.add('open');
});
document.getElementById('sidebarOverlay').addEventListener('click', () => {
    document.getElementById('sidebar').classList.remove('open');
    document.getElementById('sidebarOverlay').classList.remove('open');
});

// Sorting
document.querySelectorAll('th[data-col]').forEach(th => {
    th.addEventListener('click', () => {
        const col = th.dataset.col;
        const type = th.dataset.type;
        allTrades = Sortable.sort(allTrades, col, type);
        document.querySelectorAll('th').forEach(h => h.classList.remove('sorted'));
        th.classList.add('sorted');
        th.querySelector('.sort-icon').outerHTML = Sortable.getIcon(col);
        renderTrades();
    });
});

const loadTrades = async () => {
    const body = document.getElementById('tradesBody');
    body.innerHTML = Skeleton.rows(5);
    try {
        const data = await api.getTrades(100, 0);
        allTrades = data.trades || [];
        applyFilters();
    } catch (err) {
        Toast.error('خطا در بارگذاری تریدها');
    }
};

const applyFilters = () => {
    const status = document.getElementById('filterStatus').value;
    const side = document.getElementById('filterSide').value;
    const symbol = document.getElementById('filterSymbol').value.toUpperCase();

    let filtered = allTrades;
    if (status) filtered = filtered.filter(t => t.status === status);
    if (side) filtered = filtered.filter(t => t.side === side);
    if (symbol) filtered = filtered.filter(t => t.symbol.includes(symbol));

    allTrades = filtered;
    currentPage = 0;
    renderTrades();
};

const renderTrades = () => {
    const body = document.getElementById('tradesBody');
    const start = currentPage * PAGE_SIZE;
    const pageData = allTrades.slice(start, start + PAGE_SIZE);
    const totalPages = Math.ceil(allTrades.length / PAGE_SIZE);

    if (pageData.length === 0) {
        body.innerHTML = `<tr><td colspan="10">${EmptyState.trades()}</td></tr>`;
        document.getElementById('tradesPagination').innerHTML = '';
        return;
    }

    body.innerHTML = pageData.map(t => `
        <tr>
            <td data-label="نماد"><a href="#" onclick="openDetailModal(${t.ID});return false" style="color:var(--gold);text-decoration:underline">${t.symbol}</a></td>
            <td data-label="طرف"><span class="trade-badge ${t.side === 'LONG' ? 'long' : 'short'}">${t.side}</span></td>
            <td data-label="ورود" style="direction:ltr;text-align:right">$${t.entry_price}</td>
            <td data-label="SL" style="direction:ltr;text-align:right">$${t.stop_loss}</td>
            <td data-label="TP" style="direction:ltr;text-align:right">${t.take_profit ? '$' + t.take_profit : '—'}</td>
            <td data-label="حجم" style="direction:ltr;text-align:right">${t.quantity.toFixed(4)}</td>
            <td data-label="PnL" class="${t.pnl != null && t.pnl >= 0 ? 'up' : 'down'}" style="direction:ltr;text-align:right">
                ${t.pnl != null ? (t.pnl >= 0 ? '+' : '') + '$' + t.pnl.toFixed(2) : '—'}
            </td>
            <td data-label="وضعیت"><span class="status-badge">${t.status}</span></td>
            <td data-label="عملیات" class="action-cell">
                ${t.status === 'OPEN' ? `
                    <button class="btn btn-gold btn-sm" onclick="openCloseModal(${t.ID})">بستن</button>
                    <button class="btn btn-ghost btn-sm" onclick="openPartialModal(${t.ID})">جزئی</button>
                    <button class="btn btn-ghost btn-sm" onclick="moveBreakeven(${t.ID})">BE</button>
                    <button class="btn btn-ghost btn-sm" onclick="openEditModal(${t.ID})">ویرایش</button>
                ` : '—'}
            </td>
            <td data-label="" class="action-cell"><button class="btn btn-danger btn-sm" onclick="deleteTrade(${t.ID})"><i class="ti ti-trash"></i></button></td>
        </tr>
    `).join('');

    Pagination.render(document.getElementById('tradesPagination'), currentPage, totalPages, (page) => {
        currentPage = page;
        renderTrades();
    });
};

// Load portfolios for modals
const loadPortfolios = async () => {
    try {
        const data = await api.getPortfolios();
        const portfolios = data.portfolios || [];
        const html = portfolios.length > 0
            ? portfolios.map(p => `<option value="${p.ID}">${p.name}</option>`).join('')
            : '<option value="">ابتدا پورتفولیو بسازید</option>';
        document.getElementById('tradePortfolio').innerHTML = html;
        document.getElementById('editTradePortfolio').innerHTML = html;
    } catch {}
};

// Trade detail modal
const openDetailModal = async (id) => {
    try {
        const data = await api.getTrade(id);
        const t = data.trade;
        const pnlClass = t.pnl >= 0 ? 'up' : 'down';
        const pnlText = t.pnl != null ? (t.pnl >= 0 ? '+' : '') + '$' + t.pnl.toFixed(2) : '—';

        document.getElementById('detailTitle').textContent = t.symbol + ' - ' + t.side;
        document.getElementById('detailContent').innerHTML = `
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:16px">
                <div class="stat-card"><div class="stat-label">قیمت ورود</div><div class="stat-value" style="font-size:16px">$${t.entry_price}</div></div>
                <div class="stat-card"><div class="stat-label">قیمت خروج</div><div class="stat-value" style="font-size:16px">${t.exit_price ? '$' + t.exit_price : '—'}</div></div>
                <div class="stat-card"><div class="stat-label">حد ضرر</div><div class="stat-value" style="font-size:16px">$${t.stop_loss}</div></div>
                <div class="stat-card"><div class="stat-label">حد سود</div><div class="stat-value" style="font-size:16px">${t.take_profit ? '$' + t.take_profit : '—'}</div></div>
                <div class="stat-card"><div class="stat-label">حجم</div><div class="stat-value" style="font-size:16px">${t.quantity.toFixed(5)}</div></div>
                <div class="stat-card"><div class="stat-label">PnL</div><div class="stat-value ${pnlClass}" style="font-size:16px">${pnlText}</div></div>
                <div class="stat-card"><div class="stat-label">اهرم</div><div class="stat-value" style="font-size:16px">${t.leverage}x</div></div>
                <div class="stat-card"><div class="stat-label">وضعیت</div><div class="stat-value" style="font-size:16px"><span class="status-badge">${t.status}</span></div></div>
            </div>
            ${t.note ? `<div class="card" style="margin-bottom:12px"><div class="card-title"><i class="ti ti-note"></i>یادداشت</div><div style="font-size:13px;color:var(--text-secondary);white-space:pre-wrap">${t.note}</div></div>` : ''}
            ${t.status === 'OPEN' ? `
                <div style="display:flex;gap:8px;margin-top:12px">
                    <button class="btn btn-gold" onclick="document.getElementById('detailModalOverlay').classList.remove('open');openCloseModal(${t.ID})">بستن ترید</button>
                    <button class="btn btn-ghost" onclick="document.getElementById('detailModalOverlay').classList.remove('open');openPartialModal(${t.ID})">بستن جزئی</button>
                    <button class="btn btn-ghost" onclick="moveBreakeven(${t.ID})">انتقال به BE</button>
                </div>
            ` : ''}
        `;
        document.getElementById('detailModalOverlay').classList.add('open');

        // Load checklist for this trade
        const checklistEl = document.getElementById('detailChecklist');
        try {
            const clData = await api.getChecklist(id);
            const cl = clData.checklist;
            if (cl && cl.items) {
                let items = [];
                try { items = JSON.parse(cl.items); } catch {}
                if (items.length > 0) {
                    checklistEl.style.display = 'block';
                    checklistEl.innerHTML = `
                        <div class="card">
                            <div class="card-title"><i class="ti ti-checkbox"></i>چک‌لیست پیش از ترید ${cl.all_met ? '<span style="color:var(--green);font-size:11px">✓ تکمیل شده</span>' : ''}</div>
                            ${items.map((item, i) => `
                                <label style="display:flex;align-items:center;gap:8px;padding:6px 0;cursor:pointer;font-size:12px;color:var(--text-secondary)">
                                    <input type="checkbox" ${item.checked ? 'checked' : ''} onchange="toggleDetailChecklist(${id}, ${i}, this.checked)" style="accent-color:var(--gold);width:16px;height:16px">
                                    <span style="${item.checked ? 'text-decoration:line-through;opacity:0.5' : ''}">${item.text}</span>
                                </label>
                            `).join('')}
                        </div>
                    `;
                } else { checklistEl.style.display = 'none'; }
            } else { checklistEl.style.display = 'none'; }
        } catch { checklistEl.style.display = 'none'; }
    } catch (err) {
        Toast.error(err.message);
    }
};

// Toggle checklist item in detail modal
async function toggleDetailChecklist(tradeId, itemIndex, checked) {
    try {
        await api.updateChecklist(tradeId, itemIndex, checked);
        openDetailModal(tradeId);
    } catch (err) { Toast.error(err.message); }
}

// Checklist functions
async function loadChecklistDefaults() {
    try {
        const data = await api.getChecklistDefaults();
        checklistDefaults = data.defaults || [];
    } catch { checklistDefaults = []; }
}

function renderChecklistItems() {
    const container = document.getElementById('checklistItems');
    if (currentChecklist.length === 0) {
        container.innerHTML = '';
        document.getElementById('checklistSection').style.display = 'none';
        return;
    }
    container.innerHTML = currentChecklist.map((item, i) => `
        <label style="display:flex;align-items:center;gap:8px;padding:6px 0;cursor:pointer;font-size:12px;color:var(--text-secondary)">
            <input type="checkbox" ${item.checked ? 'checked' : ''} onchange="toggleChecklistItem(${i})" style="accent-color:var(--gold);width:16px;height:16px">
            <span style="${item.checked ? 'text-decoration:line-through;opacity:0.5' : ''}">${item.text}</span>
        </label>
    `).join('');
}

function toggleChecklistItem(index) {
    currentChecklist[index].checked = !currentChecklist[index].checked;
    renderChecklistItems();
}

function showChecklistSection() {
    const section = document.getElementById('checklistSection');
    if (checklistDefaults.length > 0) {
        section.style.display = 'block';
        currentChecklist = checklistDefaults.map(text => ({ text, checked: false }));
        renderChecklistItems();
    } else {
        section.style.display = 'none';
        currentChecklist = [];
    }
}

// New trade
document.getElementById('newTradeBtn').addEventListener('click', async () => {
    await loadChecklistDefaults();
    showChecklistSection();
    document.getElementById('modalOverlay').classList.add('open');
});
document.getElementById('modalClose').addEventListener('click', () => {
    document.getElementById('modalOverlay').classList.remove('open');
    currentChecklist = [];
});

document.getElementById('submitTrade').addEventListener('click', async () => {
    const form = document.getElementById('modalOverlay');
    Validate.clearAll(form);

    const symbol = document.getElementById('tradeSymbol').value.trim().toUpperCase();
    const side = document.getElementById('tradeSide').value;
    const entry = parseFloat(document.getElementById('tradeEntry').value);
    const sl = parseFloat(document.getElementById('tradeSL').value);
    const tp = parseFloat(document.getElementById('tradeTP').value) || null;
    const risk = parseFloat(document.getElementById('tradeRisk').value);
    const leverage = parseFloat(document.getElementById('tradeLeverage').value) || 1;
    const portfolioID = parseInt(document.getElementById('tradePortfolio').value);
    const note = document.getElementById('tradeNote').value.trim();

    let hasError = false;
    if (!symbol) { Validate.showFieldError(document.getElementById('tradeSymbol'), 'نماد الزامی است'); hasError = true; }
    if (!entry) { Validate.showFieldError(document.getElementById('tradeEntry'), 'قیمت ورود الزامی است'); hasError = true; }
    if (!sl) { Validate.showFieldError(document.getElementById('tradeSL'), 'حد ضرر الزامی است'); hasError = true; }
    if (!risk || risk <= 0) { Validate.showFieldError(document.getElementById('tradeRisk'), 'درصد ریسک الزامی است'); hasError = true; }
    if (!portfolioID) { Toast.warning('پورتفولیو را انتخاب کنید'); hasError = true; }
    if (hasError) return;

    try {
        const tradeData = { symbol, side, entry_price: entry, stop_loss: sl, take_profit: tp, risk_percent: risk, leverage, portfolio_id: portfolioID, note };
        const res = await api.createTrade(tradeData);
        const newTradeId = res.trade && res.trade.ID;

        if (newTradeId && currentChecklist.length > 0) {
            const items = currentChecklist.map(item => ({ text: item.text, checked: item.checked }));
            try { await api.createChecklist(newTradeId, items); } catch {}
        }

        document.getElementById('modalOverlay').classList.remove('open');
        currentChecklist = [];
        Toast.success('ترید ثبت شد');
        loadTrades();
    } catch (err) {
        Toast.error(err.message);
    }
});

// Close trade
const openCloseModal = (id) => {
    selectedTradeID = id;
    document.getElementById('closeModalOverlay').classList.add('open');
};
document.getElementById('closeModalClose').addEventListener('click', () => {
    document.getElementById('closeModalOverlay').classList.remove('open');
});
document.getElementById('submitClose').addEventListener('click', async () => {
    const exitPrice = parseFloat(document.getElementById('exitPrice').value);
    if (!exitPrice) { Toast.warning('قیمت خروج را وارد کنید'); return; }
    try {
        await api.closeTrade(selectedTradeID, exitPrice);
        document.getElementById('closeModalOverlay').classList.remove('open');
        Toast.success('ترید بسته شد');
        loadTrades();
    } catch (err) { Toast.error(err.message); }
});

// Partial close
const openPartialModal = (id) => {
    selectedTradeID = id;
    document.getElementById('partialModalOverlay').classList.add('open');
};
document.getElementById('partialModalClose').addEventListener('click', () => {
    document.getElementById('partialModalOverlay').classList.remove('open');
});
document.getElementById('submitPartial').addEventListener('click', async () => {
    const percent = parseFloat(document.getElementById('partialPercent').value);
    const exitPrice = parseFloat(document.getElementById('partialExitPrice').value);
    if (!percent || percent <= 0 || percent > 100) { Toast.warning('درصد معتبر وارد کنید'); return; }
    if (!exitPrice) { Toast.warning('قیمت خروج را وارد کنید'); return; }
    try {
        await api.partialClose(selectedTradeID, percent, exitPrice);
        document.getElementById('partialModalOverlay').classList.remove('open');
        Toast.success('ترید بخشی بسته شد');
        loadTrades();
    } catch (err) { Toast.error(err.message); }
});

// Move to breakeven
const moveBreakeven = async (id) => {
    const ok = await Confirm.show('انتقال به بریک‌اون', 'استاپ لاس به قیمت ورود منتقل شود؟');
    if (!ok) return;
    try {
        await api.moveToBreakeven(id);
        Toast.success('استاپ لاس به قیمت ورود منتقل شد');
        loadTrades();
    } catch (err) { Toast.error(err.message); }
};

// Edit trade
const openEditModal = async (id) => {
    selectedTradeID = id;
    try {
        const data = await api.getTrade(id);
        const t = data.trade;
        document.getElementById('editTradeSymbol').value = t.symbol;
        document.getElementById('editTradeSide').value = t.side;
        document.getElementById('editTradeEntry').value = t.entry_price;
        document.getElementById('editTradeSL').value = t.stop_loss;
        document.getElementById('editTradeTP').value = t.take_profit || '';
        document.getElementById('editTradeRisk').value = t.risk_percent;
        document.getElementById('editTradeLeverage').value = t.leverage;
        document.getElementById('editTradeNote').value = t.note || '';
        document.getElementById('editModalOverlay').classList.add('open');
    } catch (err) { Toast.error(err.message); }
};
document.getElementById('editModalClose').addEventListener('click', () => {
    document.getElementById('editModalOverlay').classList.remove('open');
});
document.getElementById('submitEdit').addEventListener('click', async () => {
    const form = document.getElementById('editModalOverlay');
    Validate.clearAll(form);

    const symbol = document.getElementById('editTradeSymbol').value.trim().toUpperCase();
    const side = document.getElementById('editTradeSide').value;
    const entry = parseFloat(document.getElementById('editTradeEntry').value);
    const sl = parseFloat(document.getElementById('editTradeSL').value);
    const tp = parseFloat(document.getElementById('editTradeTP').value) || null;
    const risk = parseFloat(document.getElementById('editTradeRisk').value);
    const leverage = parseFloat(document.getElementById('editTradeLeverage').value) || 1;
    const note = document.getElementById('editTradeNote').value.trim();

    let hasError = false;
    if (!symbol) { Validate.showFieldError(document.getElementById('editTradeSymbol'), 'نماد الزامی است'); hasError = true; }
    if (!entry) { Validate.showFieldError(document.getElementById('editTradeEntry'), 'قیمت ورود الزامی است'); hasError = true; }
    if (!sl) { Validate.showFieldError(document.getElementById('editTradeSL'), 'حد ضرر الزامی است'); hasError = true; }
    if (hasError) return;

    try {
        await api.updateTrade(selectedTradeID, { symbol, side, entry_price: entry, stop_loss: sl, take_profit: tp, risk_percent: risk, leverage, note });
        document.getElementById('editModalOverlay').classList.remove('open');
        Toast.success('ترید ویرایش شد');
        loadTrades();
    } catch (err) { Toast.error(err.message); }
});

// Delete trade
const deleteTrade = async (id) => {
    const ok = await Confirm.show('حذف ترید', 'آیا از حذف این ترید اطمینان دارید؟ این عمل قابل بازگشت نیست.');
    if (!ok) return;
    try {
        await api.deleteTrade(id);
        Toast.success('ترید حذف شد');
        loadTrades();
    } catch (err) { Toast.error(err.message); }
};

// Detail modal close
document.getElementById('detailModalClose').addEventListener('click', () => {
    document.getElementById('detailModalOverlay').classList.remove('open');
});

// Filters
document.getElementById('filterBtn').addEventListener('click', () => { loadTrades(); });
document.getElementById('filterSymbol').addEventListener('keyup', (e) => { if (e.key === 'Enter') loadTrades(); });

// Logout
document.getElementById('logoutBtn').addEventListener('click', () => api.logout());

// Keyboard shortcuts
Shortcuts.register('n', false, false, () => document.getElementById('newTradeBtn').click());
Shortcuts.init();

// Load
loadTrades();
loadPortfolios();
