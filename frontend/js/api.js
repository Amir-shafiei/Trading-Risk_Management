const BASE_URL = '/api';

(function() {
    var theme = localStorage.getItem('theme') || 'dark';
    document.documentElement.setAttribute('data-theme', theme);
})();

const getToken = () => localStorage.getItem('token');
const getRefreshToken = () => localStorage.getItem('refresh_token');
const setTokens = (accessToken, refreshToken) => {
    localStorage.setItem('token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
};
const clearTokens = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
};

let isRefreshing = false;
let refreshQueue = [];

const authHeader = () => ({
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${getToken()}`
});

const handleResponse = async (res) => {
    if (res.status === 401 && !res.url.includes('/refresh') && !res.url.includes('/login')) {
        const refreshed = await tryRefreshToken();
        if (refreshed) return null;
    }
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'خطای سرور');
    return data;
};

const tryRefreshToken = async () => {
    const refreshToken = getRefreshToken();
    if (!refreshToken) { clearTokens(); window.location.href = '/login'; return false; }
    if (isRefreshing) return false;
    isRefreshing = true;
    try {
        const res = await fetch(`${BASE_URL}/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshToken })
        });
        if (!res.ok) { clearTokens(); window.location.href = '/login'; return false; }
        const data = await res.json();
        setTokens(data.access_token, data.refresh_token);
        isRefreshing = false;
        return true;
    } catch {
        clearTokens();
        window.location.href = '/login';
        return false;
    }
};

const authFetch = async (url, options = {}) => {
    const res = await fetch(url, { ...options, headers: { ...authHeader(), ...options.headers } });
    const data = await res.json();
    if (res.status === 401 && !url.includes('/refresh')) {
        const refreshed = await tryRefreshToken();
        if (refreshed) {
            const retry = await fetch(url, { ...options, headers: { ...authHeader(), ...options.headers } });
            const retryData = await retry.json();
            if (!retry.ok) throw new Error(retryData.error || 'خطای سرور');
            return retryData;
        }
    }
    if (!res.ok) throw new Error(data.error || 'خطای سرور');
    return data;
};

const api = {

    // AUTH
    login: async (username, password) => {
        const res = await fetch(`${BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await handleResponse(res);
        if (data) setTokens(data.access_token, data.refresh_token);
        return data;
    },

    register: async (name, email, username, password) => {
        const res = await fetch(`${BASE_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, username, password })
        });
        return handleResponse(res);
    },

    logout: async () => {
        const refreshToken = getRefreshToken();
        if (refreshToken) {
            try {
                await fetch(`${BASE_URL}/logout`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ refresh_token: refreshToken })
                });
            } catch {}
        }
        clearTokens();
        window.location.href = '/login';
    },

    isLoggedIn: () => !!getToken(),

    // DASHBOARD
    getDashboard: async () => authFetch(`${BASE_URL}/dashboard`),
    getPnLHistory: async () => authFetch(`${BASE_URL}/dashboard/pnl-history`),
    getDailyPnL: async () => authFetch(`${BASE_URL}/dashboard/daily-pnl`),

    // CALCULATOR
    calculate: async (portfolioID, entryPrice, stopLoss, takeProfit, riskPercent) => {
        return authFetch(`${BASE_URL}/calculator`, {
            method: 'POST',
            body: JSON.stringify({ portfolio_id: portfolioID, entry_price: entryPrice, stop_loss: stopLoss, take_profit: takeProfit, risk_percent: riskPercent })
        });
    },

    // TRADES
    getTrades: async (limit = 10, offset = 0) => authFetch(`${BASE_URL}/trade?limit=${limit}&offset=${offset}`),
    getTrade: async (id) => authFetch(`${BASE_URL}/trade/${id}`),
    createTrade: async (data) => authFetch(`${BASE_URL}/trade`, { method: 'POST', body: JSON.stringify(data) }),
    closeTrade: async (id, exitPrice) => authFetch(`${BASE_URL}/trade/${id}/close`, { method: 'PUT', body: JSON.stringify({ exit_price: exitPrice }) }),
    partialClose: async (id, percent, exitPrice) => authFetch(`${BASE_URL}/trade/${id}/partial-close`, { method: 'PUT', body: JSON.stringify({ percent, exit_price: exitPrice }) }),
    moveToBreakeven: async (id) => authFetch(`${BASE_URL}/trade/${id}/breakeven`, { method: 'PUT' }),
    updateTrade: async (id, data) => authFetch(`${BASE_URL}/trade/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    deleteTrade: async (id) => authFetch(`${BASE_URL}/trade/${id}`, { method: 'DELETE' }),

    // PORTFOLIO
    getPortfolios: async () => authFetch(`${BASE_URL}/portfolio`),
    getPortfolio: async (id) => authFetch(`${BASE_URL}/portfolio/${id}`),
    createPortfolio: async (name, capital) => authFetch(`${BASE_URL}/portfolio`, { method: 'POST', body: JSON.stringify({ name, capital }) }),
    setDefaultPortfolio: async (id) => authFetch(`${BASE_URL}/portfolio/${id}/default`, { method: 'PUT' }),
    deletePortfolio: async (id) => authFetch(`${BASE_URL}/portfolio/${id}`, { method: 'DELETE' }),
    setDailyLossLimit: async (id, limit) => authFetch(`${BASE_URL}/portfolio/${id}/daily-loss`, { method: 'PUT', body: JSON.stringify({ max_daily_loss: limit }) }),
    setMaxOpenTrades: async (id, max) => authFetch(`${BASE_URL}/portfolio/${id}/max-open-trades`, { method: 'PUT', body: JSON.stringify({ max_open_trades: max }) }),
    getDailyLossStatus: async () => authFetch(`${BASE_URL}/portfolio/daily-loss-status`),

    // ALERTS
    checkAlerts: async () => authFetch(`${BASE_URL}/alerts/check`, { method: 'POST' }),
    getAlerts: async () => authFetch(`${BASE_URL}/alerts`),
    markAlertRead: async (id) => authFetch(`${BASE_URL}/alerts/${id}/read`, { method: 'PUT' }),

    // CHECKLIST
    createChecklist: async (tradeId, items) => authFetch(`${BASE_URL}/checklist`, { method: 'POST', body: JSON.stringify({ trade_id: tradeId, items }) }),
    getChecklist: async (tradeId) => authFetch(`${BASE_URL}/checklist/${tradeId}`),
    updateChecklist: async (tradeId, itemIndex, checked) => authFetch(`${BASE_URL}/checklist/${tradeId}`, { method: 'PUT', body: JSON.stringify({ item_index: itemIndex, checked }) }),
    getChecklistDefaults: async () => authFetch(`${BASE_URL}/checklist/defaults`),
    setChecklistDefaults: async (items) => authFetch(`${BASE_URL}/checklist/defaults`, { method: 'PUT', body: JSON.stringify({ items }) }),

    // BADGES
    checkBadges: async () => authFetch(`${BASE_URL}/badges/check`, { method: 'POST' }),
    getBadges: async () => authFetch(`${BASE_URL}/badges`),

    // NEWS
    getNews: async () => authFetch(`${BASE_URL}/news`),
    refreshNews: async () => authFetch(`${BASE_URL}/news/refresh`, { method: 'POST' }),

    // USER
    changePassword: async (oldPassword, newPassword) => {
        return authFetch(`${BASE_URL}/user/password`, {
            method: 'PUT',
            body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
        });
    }
};
