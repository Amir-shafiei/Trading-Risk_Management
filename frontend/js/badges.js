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

const BADGE_ICONS = {
    'blood': '🩸', 'fire': '🔥', 'trophy': '🏆', 'chart': '📊',
    'target': '🎯', 'money': '💰', 'shield': '🛡️', 'medal': '🏅',
    'hand': '✋', 'star': '⭐'
};

const ALL_BADGES = [
    { name: 'اولین خون', description: 'اولین ترید خود را ببندید', icon: 'blood' },
    { name: 'reaktion برنده', description: '۵ ترید برنده پشت سر هم', icon: 'fire' },
    { name: 'توقف‌ناپذیر', description: '۱۰ ترید برنده پشت سر هم', icon: 'trophy' },
    { name: 'ثبات‌قدم', description: 'نرخ برد بالای ۶۰٪ با ۲۰+ ترید', icon: 'chart' },
    { name: 'تک‌تیرانداز', description: 'نرخ برد بالای ۷۵٪ با ۱۰+ ترید', icon: 'target' },
    { name: 'برنده بزرگ', description: 'سود یک ترید بیش از ۵٪ پورتفولیو', icon: 'money' },
    { name: 'مدیر ریسک', description: '۱۰ ترید با ریسک کمتر از ۱٪', icon: 'shield' },
    { name: 'کهنه‌سرباز', description: 'تکمیل ۵۰ ترید', icon: 'medal' },
    { name: 'دست‌آهنین', description: 'نگه‌داشتن ترید بیش از ۲۴ ساعت با سود', icon: 'hand' },
    { name: 'هفته بی‌نقص', description: 'بردن همه تریدها در یک هفته (حداقل ۳)', icon: 'star' },
];

const FA_IR_MONTHS = ['فروردین', 'اردیبهشت', 'خرداد', 'تیر', 'مرداد', 'شهریور', 'مهر', 'آبان', 'آذر', 'دی', 'بهمن', 'اسفند'];

function formatFaDate(dateStr) {
    const d = new Date(dateStr);
    return d.toLocaleDateString('fa-IR', { year: 'numeric', month: 'long', day: 'numeric' });
}

async function loadBadges() {
    const grid = document.getElementById('badgesGrid');
    grid.innerHTML = Array(6).fill('<div class="skeleton skeleton-card" style="height:180px;border-radius:16px"></div>').join('');

    try {
        const data = await api.getBadges();
        const earned = data.badges || [];
        const earnedMap = {};
        earned.forEach(b => { earnedMap[b.name] = b; });

        let earnedCount = 0;
        const total = ALL_BADGES.length;

        ALL_BADGES.forEach(b => {
            if (earnedMap[b.name]) earnedCount++;
        });

        document.getElementById('badgeStats').innerHTML = `
            <div class="badge-stat-card">
                <div class="badge-stat-val">${earnedCount}</div>
                <div class="badge-stat-label">نشان کسب شده</div>
            </div>
            <div class="badge-stat-card">
                <div class="badge-stat-val">${total - earnedCount}</div>
                <div class="badge-stat-label">نشان باقی‌مانده</div>
            </div>
            <div class="badge-stat-card">
                <div class="badge-stat-val">${total > 0 ? Math.round(earnedCount / total * 100) : 0}%</div>
                <div class="badge-stat-label">پیشرفت کل</div>
            </div>
        `;

        grid.innerHTML = ALL_BADGES.map(b => {
            const isEarned = !!earnedMap[b.name];
            const earnedDate = isEarned ? formatFaDate(earnedMap[b.name].earned_at) : null;
            return `
                <div class="badge-card ${isEarned ? 'earned' : 'locked'}">
                    ${isEarned ? '<div class="badge-earned-check"><i class="ti ti-check"></i></div>' : '<div class="badge-locked-icon"><i class="ti ti-lock"></i></div>'}
                    <div class="badge-icon">${BADGE_ICONS[b.icon] || '🏅'}</div>
                    <div class="badge-name">${b.name}</div>
                    <div class="badge-desc">${b.description}</div>
                    ${isEarned ? `<div class="badge-date">${earnedDate}</div>` : ''}
                </div>
            `;
        }).join('');

    } catch (err) {
        grid.innerHTML = EmptyState.badges();
    }
}

document.getElementById('checkBadgesBtn').addEventListener('click', async () => {
    try {
        const res = await api.checkBadges();
        if (res.new_badges && res.new_badges.length > 0) {
            Toast.success(`${res.new_badges.length} نشان جدید کسب کردید! 🎉`);
        } else {
            Toast.info('نشان جدیدی کسب نشد');
        }
        loadBadges();
    } catch (err) {
        Toast.error(err.message);
    }
});

document.getElementById('logoutBtn').addEventListener('click', () => api.logout());
Shortcuts.init();

loadBadges();
