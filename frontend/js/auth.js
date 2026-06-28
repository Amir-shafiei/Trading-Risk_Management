if (api.isLoggedIn()) {
    window.location.href = '/dashboard' // ✅
}

const loginForm = document.getElementById('loginForm')
const alert = document.getElementById('alert')
const submitBtn = document.getElementById('submitBtn')

const showAlert = (msg, type) => {
    alert.textContent = msg
    alert.className = `alert ${type}`
}

loginForm.addEventListener('submit', async (e) => {
    e.preventDefault()

    const username = document.getElementById('username').value.trim()
    const password = document.getElementById('password').value.trim()

    if (!username || !password) {
        showAlert('لطفاً همه فیلدها را پر کنید', 'error')
        return
    }

    submitBtn.disabled = true
    submitBtn.textContent = 'در حال ورود...'

    try {
        await api.login(username, password)
        showAlert('ورود موفق!', 'success')
        setTimeout(() => {
            window.location.href = '/dashboard' // ✅
        }, 800)
    } catch (err) {
        showAlert(err.message, 'error')
    } finally {
        submitBtn.disabled = false
        submitBtn.textContent = 'ورود'
    }
})