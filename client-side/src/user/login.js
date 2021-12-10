import React, { Fragment, useContext, useState } from 'react'
import { Button, Row, Col, Form, FormGroup, Label, Input, CardImg } from 'reactstrap';
import { AuthContext } from '../App'
import { useNavigate } from "react-router-dom";
import './login.scss';
import axios from 'axios'
import { Container } from 'reactstrap'
import { Link } from 'react-router-dom';
const qs = require('querystring')
const api = 'http://localhost:8080'

function LoginComp(props) {
    let navigate = useNavigate();

    const initialState = {
        isSubmitting: false,
        errorMessage: null
    }

    const stateForm = {
        id: "",
        username: ""
    }
    

    const [data, setData] = useState(initialState)
    const [dataform, setDataForm] = useState(stateForm)


    const handleInputChange = event => {
        setDataForm({
            ...dataform,
            [event.target.name]: event.target.value,
        })

    }

    const handleFormSubmit = event => {
        event.preventDefault()
        
        setData({
            ...data,
            isSubmitting: true,
            errorMessage: null
        })

        const requestBody = {
            username: dataform.username
        }

        fetch('http://localhost:8080/login', {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            mode: 'cors', // no-cors, *cors, same-origin
            cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
            credentials: 'same-origin', // include, *same-origin, omit
            headers: {
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            redirect: 'follow', // manual, *follow, error
            referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
            body: JSON.stringify(requestBody)
        }).then(async response => {
            let data = await response.json();

            if ('id' in data && 'username' in data) {
                localStorage.setItem('username', data.username)
                localStorage.setItem('id', data.id)
                
                navigate("/chat");
            }
        })
    }

    return (
        <div className="login-canva">
            <Container>
                <br />
                <Row>
                    <Col>
                        <h1>Login Form</h1>

                        <hr />

                        <Form onSubmit={handleFormSubmit}>
                            <FormGroup>
                                <Label for="exampleEmail">Username</Label>
                                <Input
                                    type="username"
                                    onChange={handleInputChange}
                                    name="username"
                                    id="exampleusername"
                                    placeholder="Input Username"
                                    value={dataform.username}
                                />
                            </FormGroup>

                            <FormGroup>
                                {data.errorMessage && (
                                    <div className="alert alert-danger" role="alert">
                                        {data.errorMessage}
                                    </div>
                                )}
                            </FormGroup>
                            <FormGroup>
                                <Button disabled={data.isSubmitting}>
                                    {data.isSubmitting ? (
                                        "..Loading"
                                    ) :
                                        (
                                            "Login"
                                        )
                                    }
                                </Button>
                            </FormGroup>
                        </Form>
                    </Col>
                </Row>
            </Container>
        </div>
    )
}



export default LoginComp