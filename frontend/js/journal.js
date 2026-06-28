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

var allTrades = [];
var journalPage = 0;
var journalLimit = 20;
var hasMore = true;

async function loadJournal() {
    try {
        var data = await api.getTrades(journalLimit + 1, journalPage * journalLimit);
        var trades = data.trades || [];
        hasMore = trades.length > journalLimit;
        if (hasMore) trades = trades.slice(0, journalLimit);

        allTrades = trades.filter(function(t) { return t.status !== 'OPEN'; });
        renderJournal(allTrades);
        renderStats(allTrades);
        renderSymbolStats(allTrades);
        renderSummary(allTrades);
        renderPagination();
    } catch (err) {
        console.error(err);
    }
}

function renderJournal(trades) {
    var list = document.getElementById('journalList');
    if (!trades || trades.length === 0) {
        list.innerHTML = EmptyState.trades();
        return;
    }
    var html = '';
    for (var i = 0; i < trades.length; i++) {
        var t = trades[i];
        var pnl = t.pnl || 0;
        var isWin = pnl >= 0;
        var date = t.closed_at ? new Date(t.closed_at).toLocaleDateString('fa-IR') : '---';
        html += '<div class="journal-card ' + (isWin ? 'win' : 'loss') + '">';
        html += '<div class="journal-top"><div>';
        html += '<div style="display:flex;align-items:center;gap:8px;flex-direction:row-reverse">';
        html += '<span class="journal-sym">' + t.symbol + '</span>';
        html += '<span class="trade-badge ' + (t.side === 'LONG' ? 'long' : 'short') + '">' + t.side + '</span>';
        html += '<span class="status-badge">' + t.status + '</span>';
        html += '</div>';
        html += '<div class="journal-date">' + date + '</div>';
        html += '</div>';
        html += '<div class="journal-pnl ' + (isWin ? 'up' : 'down') + '">';
        html += (isWin ? '+' : '') + '$' + pnl.toFixed(2);
        html += '</div></div>';
        html += '<div class="journal-details">';
        html += '<div><div class="journal-detail-label">قیمت ورود</div><div class="journal-detail-val">$' + t.entry_price + '</div></div>';
        html += '<div><div class="journal-detail-label">قیمت خروج</div><div class="journal-detail-val">' + (t.exit_price ? '$' + t.exit_price : '—') + '</div></div>';
        html += '<div><div class="journal-detail-label">حجم</div><div class="journal-detail-val">' + t.quantity.toFixed(4) + '</div></div>';
        html += '<div><div class="journal-detail-label">ریسک</div><div class="journal-detail-val">' + t.risk_percent + '%</div></div>';
        html += '</div>';
        if (t.note) {
            html += '<div class="journal-note">' + t.note + '</div>';
        }
        html += '</div>';
    }
    list.innerHTML = html;
}

function renderStats(trades) {
    var pnls = [];
    for (var i = 0; i < trades.length; i++) {
        if (trades[i].pnl != null) pnls.push(trades[i].pnl);
    }
    if (pnls.length === 0) return;

    var best = Math.max.apply(null, pnls);
    var worst = Math.min.apply(null, pnls);
    var wins = pnls.filter(function(p) { return p > 0; });
    var losses = pnls.filter(function(p) { return p < 0; });
    var avgWin = 0;
    var avgLoss = 0;
    if (wins.length > 0) { avgWin = wins.reduce(function(a, b) { return a + b; }, 0) / wins.length; }
    if (losses.length > 0) { avgLoss = losses.reduce(function(a, b) { return a + b; }, 0) / losses.length; }

    document.getElementById('bestTrade').textContent = '+$' + best.toFixed(2);
    document.getElementById('worstTrade').textContent = '$' + worst.toFixed(2);
    document.getElementById('avgWin').textContent = '+$' + avgWin.toFixed(2);
    document.getElementById('avgLoss').textContent = '$' + avgLoss.toFixed(2);
}

function renderSymbolStats(trades) {
    var symbols = {};
    for (var i = 0; i < trades.length; i++) {
        var t = trades[i];
        if (!symbols[t.symbol]) { symbols[t.symbol] = { count: 0, pnl: 0 }; }
        symbols[t.symbol].count++;
        symbols[t.symbol].pnl += t.pnl || 0;
    }

    var sorted = Object.entries(symbols).sort(function(a, b) { return b[1].pnl - a[1].pnl; });
    var html = '';
    for (var i = 0; i < sorted.length; i++) {
        var sym = sorted[i][0];
        var data = sorted[i][1];
        html += '<div class="symbol-row"><div>';
        html += '<div class="symbol-name">' + sym + '</div>';
        html += '<div class="symbol-trades">' + data.count + ' ترید</div>';
        html += '</div>';
        html += '<div class="symbol-pnl ' + (data.pnl >= 0 ? 'up' : 'down') + '">';
        html += (data.pnl >= 0 ? '+' : '') + '$' + data.pnl.toFixed(2);
        html += '</div></div>';
    }
    document.getElementById('symbolStats').innerHTML = html;
}

function renderSummary(trades) {
    var pnls = [];
    for (var i = 0; i < trades.length; i++) {
        if (trades[i].pnl != null) pnls.push(trades[i].pnl);
    }
    var wins = pnls.filter(function(p) { return p > 0; });
    var losses = pnls.filter(function(p) { return p <= 0; });
    var totalPnl = pnls.reduce(function(a, b) { return a + b; }, 0);
    var winRate = pnls.length > 0 ? (wins.length / pnls.length * 100).toFixed(1) : 0;
    var avgWin = wins.length > 0 ? wins.reduce(function(a, b) { return a + b; }, 0) / wins.length : 0;
    var avgLoss = losses.length > 0 ? Math.abs(losses.reduce(function(a, b) { return a + b; }, 0) / losses.length) : 0;
    var rr = avgLoss > 0 ? (avgWin / avgLoss).toFixed(2) : '---';

    var html = '';
    html += '<div class="summary-row"><span class="summary-label">کل تریدها</span><span class="summary-val">' + trades.length + '</span></div>';
    html += '<div class="summary-row"><span class="summary-label">تریدهای سودده</span><span class="summary-val up">' + wins.length + '</span></div>';
    html += '<div class="summary-row"><span class="summary-label">تریدهای ضررده</span><span class="summary-val down">' + losses.length + '</span></div>';
    html += '<div class="summary-row"><span class="summary-label">نرخ برد</span><span class="summary-val gold">' + winRate + '%</span></div>';
    html += '<div class="summary-row"><span class="summary-label">کل PnL</span><span class="summary-val ' + (totalPnl >= 0 ? 'up' : 'down') + '">' + (totalPnl >= 0 ? '+' : '') + '$' + totalPnl.toFixed(2) + '</span></div>';
    html += '<div class="summary-row"><span class="summary-label">میانگین R/R</span><span class="summary-val gold">1 : ' + rr + '</span></div>';
    document.getElementById('summaryList').innerHTML = html;
}

document.getElementById('filterBtn').addEventListener('click', function() {
    var symbol = document.getElementById('filterSymbol').value.toUpperCase();
    var side = document.getElementById('filterSide').value;
    var result = document.getElementById('filterResult').value;
    var filtered = allTrades.slice();
    if (symbol) filtered = filtered.filter(function(t) { return t.symbol.includes(symbol); });
    if (side) filtered = filtered.filter(function(t) { return t.side === side; });
    if (result === 'win') filtered = filtered.filter(function(t) { return t.pnl > 0; });
    if (result === 'loss') filtered = filtered.filter(function(t) { return t.pnl <= 0; });
    renderJournal(filtered);
});

function renderPagination() {
    var container = document.getElementById('journalPagination');
    if (!container) return;
    var html = '';
    if (journalPage > 0) {
        html += '<button class="page-btn" onclick="changeJournalPage(' + (journalPage - 1) + ')">قبلی</button>';
    }
    html += '<span class="page-btn active">' + (journalPage + 1) + '</span>';
    if (hasMore) {
        html += '<button class="page-btn" onclick="changeJournalPage(' + (journalPage + 1) + ')">بعدی</button>';
    }
    container.innerHTML = html;
}

function changeJournalPage(page) {
    journalPage = page;
    loadJournal();
}

document.getElementById('logoutBtn').addEventListener('click', function() { api.logout(); });
Shortcuts.init();

// Heatmap
async function loadHeatmap() {
    try {
        var data = await api.getDailyPnL();
        var dailyPnL = data.daily_pnl || [];
        renderHeatmap(dailyPnL);
    } catch (err) {
        console.error(err);
    }
}

function renderHeatmap(dailyPnL) {
    var container = document.getElementById('heatmapContainer');
    if (!dailyPnL || dailyPnL.length === 0) {
        container.innerHTML = '<div class="loading">داده‌ای برای نمایش وجود ندارد</div>';
        return;
    }

    var pnlMap = {};
    var maxAbs = 0;
    for (var i = 0; i < dailyPnL.length; i++) {
        pnlMap[dailyPnL[i].date] = dailyPnL[i].pnl;
        if (Math.abs(dailyPnL[i].pnl) > maxAbs) maxAbs = Math.abs(dailyPnL[i].pnl);
    }

    var today = new Date();
    var startDate = new Date(today);
    startDate.setMonth(startDate.getMonth() - 6);
    startDate.setDate(1);

    var months = [];
    var current = new Date(startDate);
    while (current <= today) {
        var monthKey = current.getFullYear() + '-' + String(current.getMonth() + 1).padStart(2, '0');
        var dayOfWeek = current.getDay();
        var dateKey = current.getFullYear() + '-' + String(current.getMonth() + 1).padStart(2, '0') + '-' + String(current.getDate()).padStart(2, '0');
        var pnl = pnlMap[dateKey] || 0;

        if (!months[monthKey]) months[monthKey] = [];
        months[monthKey].push({ date: dateKey, pnl: pnl, dow: dayOfWeek });

        current.setDate(current.getDate() + 1);
    }

    var html = '';
    var monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

    for (var mk in months) {
        var days = months[mk];
        var parts = mk.split('-');
        var monthLabel = monthNames[parseInt(parts[1]) - 1];

        html += '<div class="heatmap-month">';
        html += '<div class="heatmap-month-label">' + monthLabel + '</div>';

        var weeks = [];
        var week = [];
        for (var d = 0; d < days.length; d++) {
            if (d === 0 && days[d].dow !== 0) {
                for (var e = 0; e < days[d].dow; e++) week.push(null);
            }
            week.push(days[d]);
            if (days[d].dow === 6 || d === days.length - 1) {
                weeks.push(week);
                week = [];
            }
        }

        for (var w = 0; w < weeks.length; w++) {
            html += '<div class="heatmap-week">';
            for (var dd = 0; dd < weeks[w].length; dd++) {
                var day = weeks[w][dd];
                if (!day) {
                    html += '<div class="heatmap-day" style="visibility:hidden"></div>';
                    continue;
                }
                var cls = 'heatmap-day';
                if (day.pnl > 0) {
                    var ratio = maxAbs > 0 ? day.pnl / maxAbs : 0;
                    if (ratio > 0.75) cls += ' p4';
                    else if (ratio > 0.5) cls += ' p3';
                    else if (ratio > 0.25) cls += ' p2';
                    else cls += ' p1';
                } else if (day.pnl < 0) {
                    var ratio = maxAbs > 0 ? Math.abs(day.pnl) / maxAbs : 0;
                    if (ratio > 0.75) cls += ' n4';
                    else if (ratio > 0.5) cls += ' n3';
                    else if (ratio > 0.25) cls += ' n2';
                    else cls += ' n1';
                }
                var tooltip = day.date + ': ' + (day.pnl >= 0 ? '+' : '') + '$' + day.pnl.toFixed(2);
                html += '<div class="' + cls + '" data-tooltip="' + tooltip + '"></div>';
            }
            html += '</div>';
        }
        html += '</div>';
    }

    container.innerHTML = html;
}

loadJournal();
loadHeatmap();