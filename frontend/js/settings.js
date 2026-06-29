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

// Password change
document.getElementById('changePasswordBtn').addEventListener('click', async () => {
    Validate.clearAll(document.querySelector('.s-card'));
    const oldPassword = document.getElementById('oldPassword').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    let hasError = false;
    if (!oldPassword) { Validate.showFieldError(document.getElementById('oldPassword'), 'رمز فعلی الزامی است'); hasError = true; }
    if (!newPassword) { Validate.showFieldError(document.getElementById('newPassword'), 'رمز جدید الزامی است'); hasError = true; }
    if (newPassword !== confirmPassword) { Validate.showFieldError(document.getElementById('confirmPassword'), 'رمزها مطابقت ندارند'); hasError = true; }
    if (newPassword && newPassword.length < 6) { Validate.showFieldError(document.getElementById('newPassword'), 'حداقل ۶ کاراکتر'); hasError = true; }
    if (hasError) return;

    try {
        await api.changePassword(oldPassword, newPassword);
        Toast.success('رمز عبور تغییر کرد');
        document.getElementById('oldPassword').value = '';
        document.getElementById('newPassword').value = '';
        document.getElementById('confirmPassword').value = '';
    } catch (err) {
        Toast.error(err.message);
    }
});

// Theme toggle
let currentTheme = localStorage.getItem('theme') || 'dark';
document.documentElement.setAttribute('data-theme', currentTheme);
updateThemeButtons(currentTheme);

document.getElementById('darkBtn').addEventListener('click', () => setTheme('dark'));
document.getElementById('lightBtn').addEventListener('click', () => setTheme('light'));

function setTheme(theme) {
    currentTheme = theme;
    localStorage.setItem('theme', theme);
    document.documentElement.setAttribute('data-theme', theme);
    updateThemeButtons(theme);
}

function updateThemeButtons(theme) {
    document.getElementById('darkBtn').classList.toggle('active', theme === 'dark');
    document.getElementById('lightBtn').classList.toggle('active', theme === 'light');
}

// Risk limits
const loadPortfolios = async () => {
    try {
        const data = await api.getPortfolios();
        const portfolios = data.portfolios || [];
        const select = document.getElementById('riskPortfolio');
        select.innerHTML = portfolios.length > 0
            ? portfolios.map(p => `<option value="${p.ID}" data-daily="${p.max_daily_loss || 0}" data-max="${p.max_open_trades || 0}">${p.name}</option>`).join('')
            : '<option value="">ابتدا پورتفولیو بسازید</option>';

        if (portfolios.length > 0) {
            selectPortfolio(portfolios[0]);
        }
    } catch {}
};

document.getElementById('riskPortfolio').addEventListener('change', (e) => {
    const opt = e.target.selectedOptions[0];
    if (opt) selectPortfolio({ max_daily_loss: parseFloat(opt.dataset.daily), max_open_trades: parseInt(opt.dataset.max) });
});

function selectPortfolio(p) {
    document.getElementById('dailyLossLimit').value = p.max_daily_loss || '';
    document.getElementById('maxOpenTrades').value = p.max_open_trades || '';
}

document.getElementById('saveRiskLimits').addEventListener('click', async () => {
    const portfolioID = parseInt(document.getElementById('riskPortfolio').value);
    if (!portfolioID) { Toast.warning('پورتفولیو را انتخاب کنید'); return; }

    const dailyLoss = parseFloat(document.getElementById('dailyLossLimit').value) || 0;
    const maxTrades = parseInt(document.getElementById('maxOpenTrades').value) || 0;

    try {
        await api.setDailyLossLimit(portfolioID, dailyLoss);
        await api.setMaxOpenTrades(portfolioID, maxTrades);
        Toast.success('محدودیت‌ها ذخیره شد');
    } catch (err) {
        Toast.error(err.message);
    }
});

document.getElementById('logoutBtn').addEventListener('click', () => api.logout());
Shortcuts.init();

loadPortfolios();

// Checklist defaults
let checklistDefaults = [];

async function loadChecklistDefaults() {
    try {
        const data = await api.getChecklistDefaults();
        checklistDefaults = data.defaults || [];
        renderChecklistDefaults();
    } catch { checklistDefaults = []; renderChecklistDefaults(); }
}

function renderChecklistDefaults() {
    const list = document.getElementById('checklistDefaultsList');
    if (checklistDefaults.length === 0) {
        list.innerHTML = '';
        return;
    }
    list.innerHTML = checklistDefaults.map((item, i) => `
        <div class="checklist-item">
            <span>${item}</span>
            <button class="s-btn-remove" onclick="removeChecklistDefault(${i})"><i class="ti ti-x"></i></button>
        </div>
    `).join('');
}

function removeChecklistDefault(index) {
    checklistDefaults.splice(index, 1);
    renderChecklistDefaults();
}

document.getElementById('addChecklistDefault').addEventListener('click', () => {
    const input = document.getElementById('newChecklistDefault');
    const text = input.value.trim();
    if (!text) return;
    checklistDefaults.push(text);
    input.value = '';
    renderChecklistDefaults();
});

document.getElementById('newChecklistDefault').addEventListener('keyup', (e) => {
    if (e.key === 'Enter') document.getElementById('addChecklistDefault').click();
});

document.getElementById('saveChecklistDefaults').addEventListener('click', async () => {
    try {
        await api.setChecklistDefaults(checklistDefaults);
        Toast.success('چک‌لیست ذخیره شد');
    } catch (err) {
        Toast.error(err.message);
    }
});

loadChecklistDefaults();
