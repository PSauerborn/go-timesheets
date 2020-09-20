
export default {
    /**
     * Function used to retrieve access token from
     * local storage. If no access token is present,
     * the client is redirected to the login page
     */
    redirectLogin: function() {
        console.log('redirecting client to login page at ' + process.env.VUE_APP_LOGIN_REDIRECT)
        window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
        return
    },
    getAccessToken: function() {
        const accessToken = localStorage.getItem('userToken')
        if (!accessToken) {
            console.log('unable to find access token in local storage. redirecting to login')
            window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
            return
        }
        return accessToken
    },
    setDevToken: function() {
        console.log('setting dev access token with value ' + process.env.VUE_APP_DEV_ACCESS_TOKEN)
        localStorage.setItem('userToken', process.env.VUE_APP_DEV_ACCESS_TOKEN)
    }
}