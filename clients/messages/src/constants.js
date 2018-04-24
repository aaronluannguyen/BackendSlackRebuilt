export const ROUTES = {
    signIn: "/",
    signUp: "/signup",
    main: "/channels/:channelName",
    generalChannel: "/channels/general"
}

export const AJAX = {
    base: "api.aaronluannguyen.me/v1/",
    signUp: AJAX.base + "users",
    signIn: AJAX.base + "users/",
    signOut: AJAX.base + "sessions/mine",
    updateFLName: AJAX.base + "users/",
    jsonApplication: "application/json"
}