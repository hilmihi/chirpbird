import React from 'react';


export class Message extends React.Component {

    render() {
        return (
            <div className='message-item'>
                <div><b>{this.props.username}</b></div>
                <span>{this.props.content}</span>
            </div>
        )
    }
}