// Toast system
const Toast = {
    container: null,
    init() {
        this.container = document.createElement('div');
        this.container.className = 'toast-container';
        document.body.appendChild(this.container);
    },
    show(message, type = 'info') {
        if (!this.container) this.init();
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.innerHTML = `<i class="ti ti-${type === 'success' ? 'check' : type === 'error' ? 'x' : type === 'warning' ? 'alert-triangle' : 'info-circle'}"></i>${message}`;
        this.container.appendChild(toast);
        setTimeout(() => toast.remove(), 3000);
    },
    success(msg) { this.show(msg, 'success'); },
    error(msg) { this.show(msg, 'error'); },
    warning(msg) { this.show(msg, 'warning'); },
    info(msg) { this.show(msg, 'info'); }
};

// Confirm dialog
const Confirm = {
    show(title, message) {
        return new Promise((resolve) => {
            const overlay = document.createElement('div');
            overlay.className = 'modal-overlay confirm-modal open';
            overlay.innerHTML = `
                <div class="modal">
                    <div class="modal-title">${title}</div>
                    <div class="confirm-text">${message}</div>
                    <div class="confirm-actions">
                        <button class="btn btn-ghost" id="confirmCancel">لغو</button>
                        <button class="btn btn-danger" id="confirmOk">تایید</button>
                    </div>
                </div>
            `;
            document.body.appendChild(overlay);
            overlay.querySelector('#confirmCancel').onclick = () => { overlay.remove(); resolve(false); };
            overlay.querySelector('#confirmOk').onclick = () => { overlay.remove(); resolve(true); };
            overlay.addEventListener('click', (e) => { if (e.target === overlay) { overlay.remove(); resolve(false); } });
        });
    }
};

// Skeleton generators
const Skeleton = {
    cards(count = 4) {
        return Array(count).fill('<div class="skeleton skeleton-card"></div>').join('');
    },
    rows(count = 5) {
        return Array(count).fill('<div class="skeleton skeleton-row"></div>').join('');
    },
    text(lines = 3) {
        const widths = ['100%', '80%', '60%', '90%', '70%'];
        return Array(lines).fill(0).map((_, i) =>
            `<div class="skeleton skeleton-text" style="width:${widths[i % widths.length]}"></div>`
        ).join('');
    }
};

// Empty state
const EmptyState = {
    trades() {
        return `<div class="empty-state">
            <div class="empty-icon"><i class="ti ti-chart-candle"></i></div>
            <div class="empty-title">هیچ تریدی ثبت نشده</div>
            <div class="empty-desc">اولین ترید خود را ثبت کنید</div>
        </div>`;
    },
    portfolios() {
        return `<div class="empty-state">
            <div class="empty-icon"><i class="ti ti-wallet"></i></div>
            <div class="empty-title">هیچ پورتفولیویی وجود ندارد</div>
            <div class="empty-desc">یک پورتفولیو جدید بسازید</div>
        </div>`;
    },
    badges() {
        return `<div class="empty-state">
            <div class="empty-icon"><i class="ti ti-medal"></i></div>
            <div class="empty-title">هنوز نشانی کسب نکرده‌اید</div>
            <div class="empty-desc">با ترید کردن نشان‌ها را آزاد کنید</div>
        </div>`;
    },
    alerts() {
        return `<div class="empty-state">
            <div class="empty-icon"><i class="ti ti-bell"></i></div>
            <div class="empty-title">اعلانی وجود ندارد</div>
            <div class="empty-desc">همه چیز عادی است</div>
        </div>`;
    },
    custom(icon, title, desc) {
        return `<div class="empty-state">
            <div class="empty-icon"><i class="ti ti-${icon}"></i></div>
            <div class="empty-title">${title}</div>
            <div class="empty-desc">${desc}</div>
        </div>`;
    }
};

// Form validation
const Validate = {
    required(value, fieldName) {
        if (!value || (typeof value === 'string' && !value.trim())) {
            return `${fieldName} الزامی است`;
        }
        return null;
    },
    number(value, fieldName, opts = {}) {
        const n = parseFloat(value);
        if (isNaN(n)) return `${fieldName} باید عدد باشد`;
        if (opts.min !== undefined && n < opts.min) return `${fieldName} باید حداقل ${opts.min} باشد`;
        if (opts.max !== undefined && n > opts.max) return `${fieldName} باید حداکثر ${opts.max} باشد`;
        return null;
    },
    showFieldError(input, message) {
        input.classList.add('error');
        let errEl = input.parentNode.querySelector('.form-error');
        if (!errEl) {
            errEl = document.createElement('div');
            errEl.className = 'form-error';
            input.parentNode.appendChild(errEl);
        }
        errEl.textContent = message;
        errEl.classList.add('show');
    },
    clearFieldError(input) {
        input.classList.remove('error');
        const errEl = input.parentNode.querySelector('.form-error');
        if (errEl) errEl.classList.remove('show');
    },
    clearAll(form) {
        form.querySelectorAll('.error').forEach(el => el.classList.remove('error'));
        form.querySelectorAll('.form-error').forEach(el => el.classList.remove('show'));
    }
};

// Pagination
const Pagination = {
    render(container, currentPage, totalPages, onPageClick) {
        if (totalPages <= 1) { container.innerHTML = ''; return; }
        let html = '';
        html += `<button class="page-btn" ${currentPage === 0 ? 'disabled' : ''} data-page="${currentPage - 1}"><i class="ti ti-chevron-right"></i></button>`;
        const start = Math.max(0, currentPage - 2);
        const end = Math.min(totalPages, start + 5);
        for (let i = start; i < end; i++) {
            html += `<button class="page-btn ${i === currentPage ? 'active' : ''}" data-page="${i}">${i + 1}</button>`;
        }
        html += `<button class="page-btn" ${currentPage >= totalPages - 1 ? 'disabled' : ''} data-page="${currentPage + 1}"><i class="ti ti-chevron-left"></i></button>`;
        container.innerHTML = html;
        container.querySelectorAll('.page-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                const page = parseInt(btn.dataset.page);
                if (!isNaN(page) && page >= 0 && page < totalPages) onPageClick(page);
            });
        });
    }
};

// Sort helper
const Sortable = {
    state: { column: null, direction: 'asc' },
    sort(data, column, type = 'string') {
        if (this.state.column === column) {
            this.state.direction = this.state.direction === 'asc' ? 'desc' : 'asc';
        } else {
            this.state.column = column;
            this.state.direction = 'asc';
        }
        const dir = this.state.direction === 'asc' ? 1 : -1;
        return [...data].sort((a, b) => {
            let va = a[column], vb = b[column];
            if (va == null) return 1;
            if (vb == null) return -1;
            if (type === 'number') return (va - vb) * dir;
            return String(va).localeCompare(String(vb)) * dir;
        });
    },
    getIcon(column) {
        if (this.state.column !== column) return '<i class="ti ti-arrows-sort sort-icon"></i>';
        return this.state.direction === 'asc'
            ? '<i class="ti ti-sort-ascending sort-icon"></i>'
            : '<i class="ti ti-sort-descending sort-icon"></i>';
    },
    reset() { this.state.column = null; this.state.direction = 'asc'; }
};

// Keyboard shortcuts
const Shortcuts = {
    handlers: {},
    register(key, ctrl, shift, callback) {
        this.handlers[`${ctrl ? 'ctrl+' : ''}${shift ? 'shift+' : ''}${key.toLowerCase()}`] = callback;
    },
    init() {
        document.addEventListener('keydown', (e) => {
            if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.tagName === 'SELECT') return;
            const key = `${e.ctrlKey ? 'ctrl+' : ''}${e.shiftKey ? 'shift+' : ''}${e.key.toLowerCase()}`;
            if (this.handlers[key]) {
                e.preventDefault();
                this.handlers[key]();
            }
            if (e.key === 'Escape') {
                document.querySelectorAll('.modal-overlay.open').forEach(m => m.classList.remove('open'));
            }
        });
    }
};

// Animated counter — animates a number from 0 to target
const Counter = {
    animate(el, target, opts = {}) {
        const duration = opts.duration || 800;
        const prefix = opts.prefix || '';
        const suffix = opts.suffix || '';
        const decimals = opts.decimals || 0;
        const start = performance.now();
        const from = 0;

        function tick(now) {
            const elapsed = now - start;
            const progress = Math.min(elapsed / duration, 1);
            const eased = 1 - Math.pow(1 - progress, 3);
            const current = from + (target - from) * eased;
            el.textContent = prefix + current.toLocaleString('en', {
                minimumFractionDigits: decimals,
                maximumFractionDigits: decimals
            }) + suffix;
            if (progress < 1) requestAnimationFrame(tick);
        }
        requestAnimationFrame(tick);
    },
    animateAll(container) {
        container.querySelectorAll('[data-count]').forEach(el => {
            const target = parseFloat(el.dataset.count);
            const prefix = el.dataset.prefix || '';
            const suffix = el.dataset.suffix || '';
            const decimals = parseInt(el.dataset.decimals) || 0;
            this.animate(el, target, { prefix, suffix, decimals });
        });
    }
};

// Button ripple effect
document.addEventListener('click', (e) => {
    const btn = e.target.closest('.btn');
    if (!btn) return;
    const rect = btn.getBoundingClientRect();
    const x = ((e.clientX - rect.left) / rect.width * 100);
    const y = ((e.clientY - rect.top) / rect.height * 100);
    btn.style.setProperty('--ripple-x', x + '%');
    btn.style.setProperty('--ripple-y', y + '%');
});

// Page transition — fade in main content on load
document.addEventListener('DOMContentLoaded', () => {
    const main = document.querySelector('.main');
    if (main) {
        main.style.opacity = '0';
        main.style.transform = 'translateY(8px)';
        requestAnimationFrame(() => {
            main.style.transition = 'opacity 0.4s ease, transform 0.4s ease';
            main.style.opacity = '1';
            main.style.transform = 'translateY(0)';
        });
    }
});

// Intersection observer for scroll animations
const ScrollReveal = {
    observer: null,
    init() {
        this.observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('revealed');
                    this.observer.unobserve(entry.target);
                }
            });
        }, { threshold: 0.1 });
    },
    observe(selector) {
        if (!this.observer) this.init();
        document.querySelectorAll(selector).forEach(el => this.observer.observe(el));
    }
};
