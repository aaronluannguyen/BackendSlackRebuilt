import React from "react";

export default class SearchedUser extends React.Component {

    render() {
        return (
            <li className="collection-item avatar">
                <img src={`${this.props.photoUrl}?s=25`} className="circle"/>
                <span className="title">{this.props.username}</span>
                <p>{`${this.props.firstName} ${this.props.lastName}`}</p>
            </li>
        )
    }
}