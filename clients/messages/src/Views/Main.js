import React from "react";
import {AJAX, ROUTES} from "../constants"
import SearchResults from "../Components/SearchResults";
import ChannelMessages from "../Components/ChannelMessages";
import NewMessage from "../Components/NewMessage";
import "./Main.css";
import Channels from "../Components/Channels";

export default class MainView extends React.Component {
    constructor(props) {
        super(props);
        this.handleSignOut = this.handleSignOut.bind(this);
        this.handleUpdate = this.handleUpdate.bind(this);
        this.handleUpdateInfo = this.handleUpdateInfo.bind(this);
        this.handleCancel = this.handleCancel.bind(this);
        this.handleUserSearchOnChange = this.handleUserSearchOnChange.bind(this);
        this.triggerUserSearch = this.triggerUserSearch.bind(this);
        this.state = {
            checkActiveSession: true,
            update: false,
            updateFirst: "",
            updateLast: "",
            userSearch: "",
            showSearch: false,
            usersResults: []
        }
    }

    componentDidMount() {
        if (this.state.checkActiveSession) {
            let url = AJAX.updateFLName + window.localStorage.getItem("id");
            fetch(url, {
                method: 'GET'
            })
                .then(res => res.json())
                .then(res => {
                    this.setState({
                        firstName: res.firstName,
                        lastName: res.lastName
                    })
                })
                .catch(err => {
                    if (err) {
                        this.props.history.push(ROUTES.signIn)
                    }
                })
                .then(() => {
                    this.setState({checkActiveSession: false})
                })
        }
    }

    handleUpdate() {
        this.setState({update: true})
    }

    handleUpdateInfo() {
        let id = window.localStorage.getItem("id");
        let url = AJAX.updateFLName + id;
        fetch(url, {
            method: 'PATCH',
            body: JSON.stringify(
                {
                    firstName: `${this.state.updateFirst}`,
                    lastName: `${this.state.updateLast}`
                }
            ),
            headers: {
                'Content-Type': AJAX.jsonApplication,
                'Authorization': window.localStorage.getItem("Authorization")
            }
        })
            .then(res => res.json())
            .then(res => {
                this.setState({
                    firstName: res.firstName,
                    lastName: res.lastName,
                    update: false,
                    updateFirst: "",
                    updateLast: ""
                });
            })
            .catch(error => {
                this.setState({updateError: error.text})
            });
    }

    handleSubmit(evt) {
        evt.preventDefault();
    }

    handleCancel() {
        this.setState({update: false})
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
                window.localStorage.removeItem("Authorization");
                this.props.history.push(ROUTES.signIn);
            }
            return res
        })
        .catch(
            err => this.setState({updateError: err.text})
        )
    }

    handleUserSearchOnChange(evt) {
        let val = evt.target.value;
        this.setState({userSearch: val});
        if (val.trim().length !== 0) {
            this.setState({showSearch: true});
            this.setState({usersResults: []});
            this.triggerUserSearch(val)
        } else {
            this.setState({showSearch: false})
        }
    }

    triggerUserSearch(val) {
        let url = AJAX.userSearch + val;
        fetch(url, {
            headers: {
                'Authorization': window.localStorage.getItem("Authorization")
            }
        })
            .then(response => response.json())
            .then(json => {
                this.setState({usersResults: json})
            })
            .catch(err => this.setState({error: err.text}))
    }

    render() {
        return (
            <div>
                <div className="container">
                    <div className="row">
                        <div className="col s12 m6">
                            <div className="card blue-grey darken-1">
                                <div className="card-content white-text">
                                    <span className="card-title">Welcome {this.state.firstName} {this.state.lastName}</span>
                                    <p>Hello there! As of version 1.0, you can update your profile or sign out!</p>
                                </div>
                                <div className="card-action">
                                    <a className="waves-effect waves-light yellow lighten-1 btn" onClick={this.handleUpdate}>Update Profile</a>
                                    <a className="waves-effect waves-light red lighten-2 btn" onClick={this.handleSignOut}>Sign Out</a>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                {
                    this.state.error ?
                        <div className="alert alert-danger">
                            {this.state.error}
                        </div> :
                        undefined
                }
                {
                    this.state.update ?
                        <div className="container">
                            <form className="col s8" onSubmit={this.handleSubmit}>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="updateFirst" type="text" className="validate"
                                               value={this.state.updateFirst}
                                               onInput={evt => this.setState({updateFirst: evt.target.value})}
                                        />
                                        <label htmlFor="email">Update First Name</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="updateLast" type="text" className="validate"
                                               value={this.state.updateLast}
                                               onInput={evt => this.setState({updateLast: evt.target.value})}
                                        />
                                        <label htmlFor="password">Update Last Name</label>
                                    </div>
                                </div>
                            </form>
                            <div>
                                {
                                    this.state.updateError ?
                                        <div className="alert alert-danger">
                                            {this.state.updateError}
                                        </div> :
                                        undefined
                                }
                                <a className="waves-effect waves-light green lighten-1 btn-large" onClick={() => this.handleUpdateInfo()}>Update</a>
                                <a className="waves-effect waves-light red lighten-2 btn-large" onClick={() => this.handleCancel()}>Cancel</a>
                            </div>
                        </div> :
                        undefined
                }
                <div className="container">
                    <div className="input-field">
                        <input id="first_name" type="text" className="validate"
                            value={this.state.userSearch}
                            onChange={evt => this.handleUserSearchOnChange(evt)}

                        />
                            <label htmlFor="first_name">Search User</label>
                    </div>
                    {
                        this.state.showSearch ?
                            <SearchResults results={this.state.usersResults}/>
                            :
                            undefined
                    }
                </div>
                <div className="container">
                    <Channels/>
                    <div id="channel-messages">
                        <ChannelMessages userInfo={this.state.userInfo} channelMessageSnap={this.state.channelMessageSnap}/>
                        <div id="bottom-channel-messages" ref={bottomPlace => {this.bottom = bottomPlace;}}/>
                    </div>
                    <div>
                        <NewMessage id="text-input" channelMessageRef={this.state.channelMessageRef} userInfo={this.state.userInfo}/>
                    </div>
                </div>
            </div>
        );
    }
}