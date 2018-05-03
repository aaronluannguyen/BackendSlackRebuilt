import React from "react";
import SearchedUser from "./SearchedUser";


export default class SearchResults extends React.Component {
    render() {
        if (!this.props.results) {
            return <p>No users found ...</p>
        }

        let users = [];
        this.props.results.forEach(user=> {
            users.push(<SearchedUser key={user.id} username={user.userName} firstName={user.firstName} userInfo={user.lastName} photoUrl={user.photoURL}/>)
        });

        return (
            <div>
                <ul className="collection">
                    {users}
                </ul>
            </div>
        )
    }
}