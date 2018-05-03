export const ROUTES = {
    signIn: "/",
    signUp: "/signup",
    main: "/channels/:channelName",
    generalChannel: "/channels/general"
}

export const AJAX = {
    base: "https://api.aaronluannguyen.me/v1/",
    signUp: "https://api.aaronluannguyen.me/v1/users",
    signIn: "https://api.aaronluannguyen.me/v1/sessions",
    signOut: "https://api.aaronluannguyen.me/v1/sessions/mine",
    updateFLName: "https://api.aaronluannguyen.me/v1/users/",
    jsonApplication: "application/json",
    userSearch: "https://api.aaronluannguyen.me/v1/users?q="
}