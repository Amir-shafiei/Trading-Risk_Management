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

// State
let allEvents = [];
let currentDate = new Date();
let activeImpact = 'all';
let activeCurrency = 'all';
let autoRefreshInterval = null;

// Persian translations for common economic events
const EVENT_INFO = {
    'Non-Farm Payrolls': {
        fa: 'اشغال غیرکشاورزی (NFP)',
        desc: 'تعداد شاغلین اضافه شده در بخش غیرکشاورزی اقتصاد آمریکا. مهم‌ترین شاخص بازار کار آمریکا است. اگر عدد واقعی بالاتر از پیش‌بینی باشد، دلار تقویت می‌شود و معمولاً بیت‌کوین و طلا کاهش می‌یابند. اگر پایین‌تر باشد، دلار ضعیف و کریپتو تقویت می‌شود.'
    },
    'CPI': {
        fa: 'شاخص قیمت مصرف‌کننده (CPI)',
        desc: 'شاخص تورم آمریکا که تغییرات قیمت کالاها و خدمات را اندازه‌گیری می‌کند. اگر CPI بالاتر از پیش‌بینی باشد، فدرال رزرو احتمالاً نرخ بهره را افزایش می‌دهد که برای کریپتو منفی است. اگر پایین‌تر باشد، انتظار کاهش نرخ بهره و رشد بازارها.'
    },
    'Core CPI': {
        fa: 'شاخص قیمت مصرف‌کننده هسته‌ای',
        desc: 'نسخه CPI بدون احتساب غذا و انرژی که نماگر دقیق‌تر تورم پایدار است. فدرال رزرو بیشتر روی این شاخص تمرکز دارد. عدد بالاتر = احتمال افزایش نرخ بهره = فشار نزولی بر بیت‌کوین.'
    },
    'FOMC Statement': {
        fa: 'بیانیه کمیته فدرال رزرو',
        desc: 'بیانیه رسمی فدرال رزرو درباره سیاست پولی. بازارها به شدت به این بیانیه واکنش نشان می‌دهند. اگر لحن سخت‌گیرانه (hawkish) باشد، دلار تقویت و کریپتو ضعیف می‌شود. اگر لحن ملایم (dovish) باشد، بازار رشد می‌کند.'
    },
    'Federal Funds Rate': {
        fa: 'نرخ بهره فدرال رزرو',
        desc: 'نرخ بهره اصلی فدرال رزرو. افزایش نرخ بهره = دلار قوی‌تر و فشار بر کریپتو. کاهش نرخ = بازار صعودی. مهم‌ترین عامل قیمت بیت‌کوین در کوتاه‌مدت.'
    },
    'Interest Rate Decision': {
        fa: 'تصمیم نرخ بهره',
        desc: 'تصمیم بانک مرکزی درباره نرخ بهره. هر بانک مرکزی (ECB, BOE, BOJ و...) نرخ بهره خود را تعیین می‌کند. افزایش = ارز آن کشور تقویت. کاهش = ارز ضعیف‌تر.'
    },
    'GDP': {
        fa: 'تولید ناخالص داخلی',
        desc: 'شاخص رشد اقتصادی. عدد بالاتر از پیش‌بینی نشان‌دهنده اقتصاد قوی و معمولاً تقویت ارز آن کشور است. عدد پایین‌تر = رکود احتمالی.'
    },
    'Unemployment Rate': {
        fa: 'نرخ بیکاری',
        desc: 'درصد بیکاران جمعیت فعال. نرخ بیکاری پایین‌تر = اقتصاد قوی = تقویت دلار. نرخ بالاتر = اقتصاد ضعیف = احتمال کاهش نرخ بهره.'
    },
    'Retail Sales': {
        fa: 'فروش خرده‌فروشی',
        desc: 'میزان خرید مصرف‌کنندگان. عدد بالاتر = مصرف قوی = اقتصاد سالم = تقویت ارز. عدد پایین‌تر = کاهش مصرف.'
    },
    'PMI': {
        fa: 'شاخص مدیران خرید',
        desc: 'شاخص سلامت بخش تولید و خدمات. عدد بالای 50 = رشد اقتصادی. عدد زیر 50 = رکود. بسیار تأثیرگذار بر بازار فارکس.'
    },
    'Manufacturing PMI': {
        fa: 'شاخص مدیران خرید بخش تولید',
        desc: 'شاخص سلامت بخش تولید. بالای 50 = توسعه، زیر 50 = انقباض. تأثیر مستقیم بر ارز آن کشور.'
    },
    'Services PMI': {
        fa: 'شاخص مدیران خرید بخش خدمات',
        desc: 'شاخص سلامت بخش خدمات. بزرگ‌ترین بخش اقتصاد آمریکا. بالای 50 = رشد.'
    },
    'Consumer Confidence': {
        fa: 'اعتماد مصرف‌کننده',
        desc: 'میزان اعتماد مردم به اقتصاد. عدد بالاتر = مصرف بیشتر = اقتصاد قوی. پیش‌بینی‌کننده خوب رشد آینده.'
    },
    'Trade Balance': {
        fa: 'تراز تجاری',
        desc: 'تفاوت صادرات و واردات. مازاد تجاری = تقویت ارز. کسری تجاری = فشار بر ارز.'
    },
    'Initial Jobless Claims': {
        fa: 'درخواست اولیه بیمه بیکاری',
        desc: 'تعداد افرادی که برای اولین بار بیمه بیکاری دریافت کرده‌اند. عدد پایین‌تر = بازار کار قوی.'
    },
    'Core Retail Sales': {
        fa: 'فروش خرده‌فروشی هسته‌ای',
        desc: 'فروش خرده‌فروشی بدون احتساب خودرو. شاخص دقیق‌تر مصرف.'
    },
    'Average Earnings': {
        fa: 'میانگین درآمد',
        desc: 'تغییرات میانگین درآمد کارگران. افزایش = فشار تورمی = احتمال افزایش نرخ بهره.'
    },
    'BOE Interest Rate': {
        fa: 'نرخ بهره بانک انگلیس',
        desc: 'تصمیم بانک مرکزی انگلیس درباره نرخ بهره. تأثیر مستقیم بر پوند (GBP).'
    },
    'ECB Interest Rate': {
        fa: 'نرخ بهره بانک مرکزی اروپا',
        desc: 'تصمیم ECB درباره نرخ بهره. تأثیر مستقیم بر یورو (EUR).'
    },
    'BOJ Interest Rate': {
        fa: 'نرخ بهره بانک ژاپن',
        desc: 'تصمیم بانک مرکزی ژاپن. تأثیر مستقیم بر ین (JPY). تغییرات غیرمنتظره می‌تواند بازار را تکان دهد.'
    },
    'RBA Interest Rate': {
        fa: 'نرخ بهره بانک استرالیا',
        desc: 'تصمیم بانک مرکزی استرالیا. تأثیر بر دلار استرالیا (AUD).'
    },
    'BOC Interest Rate': {
        fa: 'نرخ بهره بانک کانادا',
        desc: 'تصمیم بانک مرکزی کانادا. تأثیر بر دلار کانادا (CAD).'
    },
    'CPI (YoY)': {
        fa: 'شاخص قیمت مصرف‌کننده (سالانه)',
        desc: 'تغییرات سالانه قیمت‌ها. شاخص اصلی تورم. بالاتر = فشار تورمی.'
    },
    'CPI (MoM)': {
        fa: 'شاخص قیمت مصرف‌کننده (ماهانه)',
        desc: 'تغییرات ماهانه قیمت‌ها. حساس‌تر به تغییرات کوتاه‌مدت.'
    },
    'Core CPI (MoM)': {
        fa: 'شاخص قیمت مصرف‌کننده هسته‌ای (ماهانه)',
        desc: 'تورم بدون غذا و انرژی. محبوب‌ترین شاخص فدرال رزرو.'
    },
    'ECB Press Conference': {
        fa: 'کنفرانس خبری بانک مرکزی اروپا',
        desc: 'رئیس ECB درباره سیاست پولی صحبت می‌کند. لحن و کلمات انتخابی تأثیر زیادی بر یورو دارد.'
    },
    'Fed Chair Speaks': {
        fa: 'سخنرانی رئیس فدرال رزرو',
        desc: 'اظهارات رئیس فدرال رزرو. هر کلمه با دقت بررسی می‌شود. لحن hawkish/dovish تعیین‌کننده جهت بازار.'
    },
    'Building Permits': {
        fa: 'مجوزهای ساخت‌وساز',
        desc: 'تعداد مجوزهای صادر شده برای ساخت خانه. شاخص پیشرو اقتصاد مسکن.'
    },
    'Housing Starts': {
        fa: 'شروع ساخت مسکن',
        desc: 'تعداد واحدهای مسکن جدید در حال ساخت. شاخص سلامت بخش مسکن.'
    },
    'Industrial Production': {
        fa: 'تولید صنعتی',
        desc: 'میزان تولید کارخانه‌ها و معادن. افزایش = رشد اقتصادی.'
    },
    'Michigan Consumer Sentiment': {
        fa: 'اعتماد مصرف‌کننده میشیگان',
        desc: 'شاخص اعتماد مصرف‌کننده آمریکا. پیش‌بینی‌کننده رفتار مصرفی آینده.'
    },
    'ADP Non-Farm Employment': {
        fa: 'اشغال ADP (غیرکشاورزی)',
        desc: 'گزارش خصوصی اشتغال آمریکا. پیش‌نمایش NFP رسمی.'
    },
    'Empire State Manufacturing': {
        fa: 'شاخص تولید امپایر استیت',
        desc: 'شاخص تولید منطقه نیویورک. بالای صفر = رشد.'
    },
    'Philly Fed Manufacturing': {
        fa: 'شاخص تولید فیلادلفیا',
        desc: 'شاخص تولید منطقه فیلادلفیا. بالای صفر = رشد بخش تولید.'
    },
    'CPI': {
        fa: 'شاخص قیمت مصرف‌کننده',
        desc: 'شاخص اصلی تورم آمریکا. تأثیر مستقیم بر سیاست فدرال رزرو و قیمت بیت‌کوین.'
    },
    'GDP (QoQ)': {
        fa: 'تولید ناخالص داخلی (فصلی)',
        desc: 'نرخ رشد اقتصادی فصلی. بالای 2٪ = رشد سالم.'
    },
    'Trade Balance': {
        fa: 'تراز تجاری',
        desc: 'مازاد یا کسری تجاری. کسری = واردات بیشتر از صادرات.'
    },
    'Federal Budget Balance': {
        fa: 'تراز بودجه فدرال',
        desc: 'مازاد یا کسری بودجه دولت آمریکا.'
    },
    'Labor Cost Index': {
        fa: 'شاخص هزینه نیروی کار',
        desc: 'تغییرات هزینه نیروی کار. افزایش = فشار تورمی.'
    },
    'Existing Home Sales': {
        fa: 'فروش خانه‌های موجود',
        desc: 'تعداد خانه‌های فروخته شده. شاخص بازار مسکن.'
    },
    'New Home Sales': {
        fa: 'فروش خانه‌های جدید',
        desc: 'تعداد خانه‌های جدید فروخته شده. شاخص پیشرو اقتصاد.'
    },
    'Durable Goods Orders': {
        fa: 'سفارشات کالاهای بادوام',
        desc: 'سفارشات کالاهایی با عمر بیش از 3 سال. شاخص سرمایه‌گذاری تجاری.'
    },
    'JOLTS Job Openings': {
        fa: 'موقعیت‌های شغلی خالی',
        desc: 'تعداد موقعیت‌های شغلی خالی. بالاتر = بازار کار رقابتی.'
    },
    ' Challenger Job Cuts': {
        fa: 'اخراج‌های اعلام شده',
        desc: 'تعداد اخراج‌های اعلام شده توسط شرکت‌ها. افزایش = ضعف بازار کار.'
    }
};

// Get event info (English + Persian)
function getEventInfo(eventName) {
    for (const [key, info] of Object.entries(EVENT_INFO)) {
        if (eventName.toLowerCase().includes(key.toLowerCase())) {
            return info;
        }
    }
    return { fa: '', desc: '' };
}

// Persian day names
const WEEKDAYS_FA = ['یکشنبه', 'دوشنبه', 'سه‌شنبه', 'چهارشنبه', 'پنجشنبه', 'جمعه', 'شنبه'];
const MONTHS_FA = ['ژانویه', 'فوریه', 'مارس', 'آوریل', 'مه', 'ژوئن', 'ژوئیه', 'اوت', 'سپتامبر', 'اکتبر', 'نوامبر', 'دسامبر'];

function formatDateDisplay(d) {
    return `${d.getFullYear()}/${String(d.getMonth()+1).padStart(2,'0')}/${String(d.getDate()).padStart(2,'0')}`;
}

function formatDateKey(d) {
    return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`;
}

function formatWeekday(d) {
    return WEEKDAYS_FA[d.getDay()];
}

// Update day navigation display
function updateDayNav() {
    document.getElementById('currentDay').textContent = formatDateDisplay(currentDate);
    document.getElementById('currentWeekday').textContent = formatWeekday(currentDate);

    const today = new Date();
    today.setHours(0,0,0,0);
    const check = new Date(currentDate);
    check.setHours(0,0,0,0);
    document.getElementById('todayBtn').style.display = today.getTime() === check.getTime() ? 'none' : 'inline-block';
}

// Render calendar for current day
function renderCalendar() {
    const body = document.getElementById('calendarBody');
    const dateKey = formatDateKey(currentDate);

    let dayEvents = allEvents.filter(e => e.date === dateKey);

    if (activeImpact !== 'all') {
        dayEvents = dayEvents.filter(e => e.impact_level === parseInt(activeImpact));
    }
    if (activeCurrency !== 'all') {
        dayEvents = dayEvents.filter(e => e.currency.toUpperCase() === activeCurrency);
    }

    if (dayEvents.length === 0) {
        body.innerHTML = `<tr><td colspan="7">
            <div class="no-events">
                <div class="no-events-icon">📅</div>
                <div>رویداد اقتصادی برای این روز وجود ندارد</div>
            </div>
        </td></tr>`;
        return;
    }

    body.innerHTML = dayEvents.map(e => {
        const impactClass = e.impact_level === 3 ? 'high' : e.impact_level === 2 ? 'medium' : 'low';
        const currClass = e.currency.replace('/', '').substring(0, 3).toUpperCase();
        const info = getEventInfo(e.event_name);

        return `<tr>
            <td data-label="ساعت" class="event-time">${e.time || '—'}</td>
            <td data-label="ارز"><span class="event-currency ${currClass}">${e.currency}</span></td>
            <td data-label="اهمیت"><span class="impact-badge ${impactClass}"><span class="impact-dot ${impactClass}"></span></span></td>
            <td data-label="رویداد" class="event-name-cell" onclick="openEventDetail(this)" data-event='${JSON.stringify(e).replace(/'/g, "&#39;")}'>
                <span class="event-name-en">${e.event_name}</span>
                ${info.fa ? `<span class="event-name-fa">${info.fa}</span>` : ''}
            </td>
            <td data-label="پیش‌بینی" class="event-value">${e.forecast || '—'}</td>
            <td data-label="قبلی" class="event-value">${e.previous || '—'}</td>
            <td data-label="واقعی" class="event-value event-actual ${getActualClass(e.actual, e.forecast)}">${e.actual || '—'}</td>
        </tr>`;
    }).join('');
}

function getActualClass(actual, forecast) {
    if (!actual || !forecast) return '';
    const a = parseFloat(actual.replace(/[^0-9.-]/g, ''));
    const f = parseFloat(forecast.replace(/[^0-9.-]/g, ''));
    if (isNaN(a) || isNaN(f)) return '';
    return a > f ? 'positive' : a < f ? 'negative' : '';
}

// Open event detail modal
function openEventDetail(td) {
    const e = JSON.parse(td.dataset.event);
    const info = getEventInfo(e.event_name);
    const impactClass = e.impact_level === 3 ? 'high' : e.impact_level === 2 ? 'medium' : 'low';
    const impactText = e.impact_level === 3 ? 'بالا' : e.impact_level === 2 ? 'متوسط' : 'پایین';

    document.getElementById('eventModalContent').innerHTML = `
        <div class="event-detail-header">
            <div>
                <div class="event-detail-title">${e.event_name}</div>
                ${info.fa ? `<div class="event-detail-title-fa">${info.fa}</div>` : ''}
            </div>
            <span class="impact-badge ${impactClass}"><span class="impact-dot ${impactClass}"></span>${impactText}</span>
        </div>
        <div class="event-detail-meta">
            <div class="event-detail-box">
                <div class="event-detail-box-label">پیش‌بینی</div>
                <div class="event-detail-box-value">${e.forecast || '—'}</div>
            </div>
            <div class="event-detail-box">
                <div class="event-detail-box-label">قبلی</div>
                <div class="event-detail-box-value">${e.previous || '—'}</div>
            </div>
            <div class="event-detail-box">
                <div class="event-detail-box-label">واقعی</div>
                <div class="event-detail-box-value event-actual ${getActualClass(e.actual, e.forecast)}">${e.actual || ' منتشر نشده'}</div>
            </div>
        </div>
        ${info.desc ? `<div class="event-detail-desc"><strong>توضیح:</strong><br>${info.desc}</div>` : '<div class="event-detail-desc" style="opacity:0.5">توضیحات فارسی برای این رویداد موجود نیست.</div>'}
    `;
    document.getElementById('eventModal').classList.add('open');
}

document.getElementById('eventModalClose').addEventListener('click', () => {
    document.getElementById('eventModal').classList.remove('open');
});

// Day navigation
document.getElementById('prevDay').addEventListener('click', () => {
    currentDate.setDate(currentDate.getDate() - 1);
    updateDayNav();
    renderCalendar();
});

document.getElementById('nextDay').addEventListener('click', () => {
    currentDate.setDate(currentDate.getDate() + 1);
    updateDayNav();
    renderCalendar();
});

document.getElementById('todayBtn').addEventListener('click', () => {
    currentDate = new Date();
    updateDayNav();
    renderCalendar();
});

// Impact filter
document.querySelectorAll('.filter-chip[data-impact]').forEach(chip => {
    chip.addEventListener('click', () => {
        document.querySelectorAll('.filter-chip[data-impact]').forEach(c => c.classList.remove('active'));
        chip.classList.add('active');
        activeImpact = chip.dataset.impact;
        renderCalendar();
    });
});

// Currency filters
function buildCurrencyFilters() {
    const currencies = [...new Set(allEvents.map(e => e.currency.toUpperCase()))].sort();
    const container = document.getElementById('currencyFilters');
    container.innerHTML = currencies.map(c =>
        `<button class="currency-btn" data-currency="${c}">${c}</button>`
    ).join('');

    container.querySelectorAll('.currency-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            container.querySelectorAll('.currency-btn').forEach(b => b.classList.remove('active'));
            if (activeCurrency === btn.dataset.currency) {
                activeCurrency = 'all';
            } else {
                btn.classList.add('active');
                activeCurrency = btn.dataset.currency;
            }
            renderCalendar();
        });
    });
}

// Load calendar data
async function loadCalendar() {
    const body = document.getElementById('calendarBody');
    body.innerHTML = `<tr><td colspan="7">${Skeleton.rows(5)}</td></tr>`;

    try {
        const data = await api.getNews();
        allEvents = data.events || [];
        buildCurrencyFilters();
        renderCalendar();
    } catch (err) {
        body.innerHTML = `<tr><td colspan="7" class="loading">خطا در دریافت تقویم</td></tr>`;
    }
}

// Auto-refresh every 5 minutes for live actual values
function startAutoRefresh() {
    if (autoRefreshInterval) clearInterval(autoRefreshInterval);
    autoRefreshInterval = setInterval(async () => {
        try {
            await api.refreshNews();
            const data = await api.getNews();
            allEvents = data.events || [];
            buildCurrencyFilters();
            renderCalendar();
        } catch {}
    }, 5 * 60 * 1000);
}

// Manual refresh
document.getElementById('refreshBtn').addEventListener('click', async () => {
    try {
        Toast.info('تقویم در حال بروزرسانی...');
        await api.refreshNews();
        const data = await api.getNews();
        allEvents = data.events || [];
        buildCurrencyFilters();
        renderCalendar();
        Toast.success('تقویم بروزرسانی شد');
    } catch (err) {
        Toast.error(err.message);
    }
});

document.getElementById('logoutBtn').addEventListener('click', () => api.logout());
Shortcuts.init();

// Init
updateDayNav();
loadCalendar();
startAutoRefresh();
