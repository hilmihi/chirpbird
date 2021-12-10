import React, { useState } from 'react';
import { Channel } from './Channel';

export class ChannelList extends React.Component {
    state = { showModal: false };

    newMessage = () => {
        this.props.onShowModal();
    }

    handleClick = id => {
        this.props.onSelectChannel(id);
    }

    render() {

        let list = <div className="no-content-message">There is no channels to show</div>;
        if (this.props.channels && this.props.channels.map) {
            list = this.props.channels.map(c => <Channel key={c.room} id={c.id} name={c.Public.rn} participants={0} onClick={this.handleClick} />);
        }
        return (
            <div className='channel-list'>
                <div className='list'>
                    {list}
                </div>
                {this.props.channels &&
                    <div className="add-chat">
                        <button onClick={this.newMessage}>New Message Room</button>
                    </div>
                }
            </div>);
    }

}