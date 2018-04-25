import React from "react";
import {Link} from "react-router-dom";
import {AJAX, ROUTES} from "../constants"

export default class MainView extends React.Component {
    constructor(props) {
        super(props);
        this.handleSignOut = this.handleSignOut.bind(this);
        this.state = {
            firstName: window.localStorage.getItem("firstName"),
            lastName: window.localStorage.getItem("lastName")
        }
    }

    componentWillMount() {
        let url = AJAX.updateFLName + window.localStorage.getItem("id");
        fetch(url, {
            method: 'GET'
        })
        .then(res => res.json())
        .catch(err => {
            if (err) {
                this.props.history.push(ROUTES.signIn)
            }
        })
    }

    handleUpdate() {

    }

    handleSignOut() {
        fetch(AJAX.signOut , {
                method: 'DELETE',
                headers: {
                    'Authorization': window.localStorage.getItem("Authorization")
                }
            }
        )
        .then((res) => {
            if (res.status !== 403) {
                window.localStorage.removeItem("id");
                window.localStorage.removeItem("username");
                window.localStorage.removeItem("firstName");
                window.localStorage.removeItem("lastName");
                window.localStorage.removeItem("Authorization");
                this.props.history.push(ROUTES.signIn);
            }
            return res
        })
        .catch(
            err => console.log(err)
        )
    }

    render() {
        return (
            <div className="container">
                <div className="row">
                    <div className="col s12 m6">
                        <div className="card blue-grey darken-1">
                            <div className="card-content white-text">
                                <span className="card-title">Welcome {this.state.firstName} {this.state.lastName}</span>
                                <p>Hello there! As of version 1.0, you can update your profile or sign out!</p>
                            </div>
                            <div className="card-action">
                                <a className="waves-effect waves-light yellow lighten-1 btn">Update Profile</a>
                                <a className="waves-effect waves-light red lighten-2 btn" onClick={this.handleSignOut}>Sign Out</a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}