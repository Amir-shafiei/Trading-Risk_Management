if (api.isLoggedIn()) {
    window.location.href = '/dashboard' // ✅
}

const registerForm = document.getElementById('registerForm')
const alertEl = document.getElementById('alert')
const submitBtn = document.getElementById('submitBtn')

const showAlert = (msg, type) => {
    alertEl.textContent = msg
    alertEl.className = `alert ${type}`
}

const validateForm = (name, email, username, password, confirmPassword) => {
    if (!name || !email || !username || !password || !confirmPassword) {
        showAlert('لطفاً همه فیلدها را پر کنید', 'error')
        return false
    }
    if (password.length < 8) {
        showAlert('رمز عبور باید حداقل ۸ کاراکتر باشد', 'error')
        return false
    }
    if (password !== confirmPassword) {
        showAlert('رمز عبور و تکرار آن یکسان نیستند', 'error')
        return false
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email)) {
        showAlert('ایمیل معتبر نیست', 'error')
        return false
    }
    return true
}

registerForm.addEventListener('submit', async (e) => {
    e.preventDefault()

    const name = document.getElementById('name').value.trim()
    const email = document.getElementById('email').value.trim()
    const username = document.getElementById('username').value.trim()
    const password = document.getElementById('password').value.trim()
    const confirmPassword = document.getElementById('confirmPassword').value.trim()

    if (!validateForm(name, email, username, password, confirmPassword)) return

    submitBtn.disabled = true
    submitBtn.textContent = 'در حال ثبت‌نام...'

    try {
        await api.register(name, email, username, password)
        showAlert('ثبت‌نام موفق! در حال انتقال...', 'success')
        setTimeout(() => {
            window.location.href = '/login' // ✅
        }, 1000)
    } catch (err) {
        showAlert(err.message, 'error')
    } finally {
        submitBtn.disabled = false
        submitBtn.textContent = 'ثبت‌نام'
    }
})