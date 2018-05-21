import React from "react";
import {Link} from "react-router-dom";
import {ROUTES} from "../constants"
import {AJAX} from "../constants";

export default class SignInView extends React.Component {
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

    handleSignUn() {
        fetch(`${AJAX.signUp}`, {
                method: 'POST',
                body: JSON.stringify(
                    {
                        email: `${this.state.userEmail}`,
                        password: `${this.state.userPassword}`,
                        passwordConf: `${this.state.userPasswordConf}`,
                        userName: `${this.state.username}`,
                        firstName: `${this.state.firstName}`,
                        lastName: `${this.state.lastName}`
                    }
                ),
                headers: {
                    'Content-Type': `${AJAX.jsonApplication}`
                }
            }
        ).then(
            (res) => {
                if (res.status < 300) {
                    let authContent = res.headers.get("Authorization");
                    localStorage.setItem("Authorization", authContent);
                    this.props.history.push(ROUTES.generalChannel);
                    return res.json()
                } else {
                    throw res
                }
            })
            .then(resJson => {
                localStorage.setItem("id", resJson.id);
            })
            .catch(error => {
                error.text().then(errMsg => {
                    this.setState({error: errMsg})
                })
            });
    }

    handleSubmit(evt) {
        evt.preventDefault();
    }

    render() {
        return (
            <div className="row">
                <div className="col s12">
                    <div id="form-container" className="container">
                        <div className="row">
                            <form className="col s8" onSubmit={this.handleSubmit}>
                                <div className="row">
                                    <div className="input-field col s6">
                                        <input id="firstName" type="text" className="validate"
                                               value={this.state.firstName}
                                               onInput={evt => this.setState({firstName: evt.target.value})}
                                        />
                                        <label htmlFor="text">First Name</label>
                                    </div>
                                    <div className="input-field col s6">
                                        <input id="lastName" type="text" className="validate"
                                               value={this.state.lastName}
                                               onInput={evt => this.setState({lastName: evt.target.value})}
                                        />
                                        <label htmlFor="text">Last Name</label>
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
                                        <label htmlFor="password">Username</label>
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
                                        <input id="passwordConf" type="password" className="validate"
                                               value={this.state.userPasswordConf}
                                               onInput={evt => this.setState({userPasswordConf: evt.target.value})}
                                        />
                                        <label htmlFor="password">Password Confirm</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        {
                            this.state.error ?
                                <div className="alert alert-danger">
                                    {this.state.error}
                                </div> :
                                undefined
                        }
                        <div>
                            <a className="waves-effect waves-light btn-large" onClick={() => this.handleSignUn()}>Sign Up</a>
                            Don't have an account yet? <Link to={ROUTES.signIn}> Sign In </Link>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}