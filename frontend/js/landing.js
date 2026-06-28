// اگه لاگین بود برو داشبورد
if (typeof api !== 'undefined' && api.isLoggedIn()) {
    window.location.href = '/dashboard'
}

// هدر scroll effect
window.addEventListener('scroll', () => {
    const header = document.querySelector('.header')
    if (window.scrollY > 20) {
        header.classList.add('scrolled')
    } else {
        header.classList.remove('scrolled')
    }
})