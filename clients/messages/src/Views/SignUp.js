import React from "react";
import {Link} from "react-router-dom";
import {AJAX, ROUTES} from "../constants"

export default class SignUpView extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            userEmail: "",
            userPassword: "",
            userPasswordConf: "",
            username: "",
            firstName: "",
            lastName: ""
        }
    }

    handleSignUp(evt) {
        evt.preventDefault();

        fetch(
            `${AJAX.signIn}`,
            {
                method: "POST",
                headers: {
                    "Content-Type": `${AJAX.jsonApplication}`
                },
                body: JSON.stringify({
                    email: `${this.state.userEmail}`,
                    password: `${this.state.userPassword}`,
                    passwordConf: `${this.state.userPasswordConf}`,
                    userName: `${this.state.username}`,
                    firstName: `${this.state.firstName}`,
                    lastName: `${this.state.lastName}`
                })
            }
        ).then(res => {

        })
            .then(
                (result) => {

                },
                (error) => {

                }
            )

        this.props.history.push(ROUTES.generalChannel);
    }

    render() {
        return (
            <div className="row">
                <div className="col s4">

                </div>
                <div className="col s8">
                    <div id="form-container" className="container">
                        <div className="row">
                            <form className="col s12" onSubmit={(evt) => this.handleSignUp(evt)}>
                                <div className="row">
                                    <div className="input-field col s6">
                                        <input id="first_name" type="text" className="validate"
                                           value={this.state.firstName}
                                           onInput={evt => this.setState({firstName: evt.target.value})}
                                        />
                                        <label htmlFor="first_name">First Name</label>
                                    </div>
                                    <div className="input-field col s6">
                                        <input id="last_name" type="text" className="validate"
                                           value={this.state.lastName}
                                           onInput={evt => this.setState({lastName: evt.target.value})}
                                        />
                                        <label htmlFor="last_name">Last Name</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="email" type="email" className="validate"
                                           value={this.state.userEmail}
                                           onInput={evt => this.setState({userEmail: evt.target.value})}
                                        />
                                        <label htmlFor="email">Email</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="username" type="text" className="validate"
                                           value={this.state.username}
                                           onInput={evt => this.setState({username: evt.target.value})}
                                        />
                                        <label htmlFor="username">Username</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"
                                           value={this.state.userPassword}
                                           onInput={evt => this.setState({userPassword: evt.target.value})}
                                        />
                                        <label htmlFor="password">Password</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"
                                           value={this.state.userPasswordConf}
                                           onInput={evt => this.setState({email: evt.target.value})}
                                        />
                                        <label htmlFor="password">Confirm Password</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div>
                            <div>
                                <a className="waves-effect waves-light btn-large">Sign Up</a>
                            </div>
                            <div>
                                Already have an account? <Link to={ROUTES.signIn}> Sign In </Link>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}