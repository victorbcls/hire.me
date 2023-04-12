import React, { useState } from "react";
import axios from "axios";
import styled from "styled-components";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
  padding: 30px;
`;

const Input = styled.input`
  width: 100%;
  padding: 10px;
  border-radius: 5px;
  border: 1px solid #ccc;
  margin-bottom: 20px;
`;

const Button = styled.button`
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  background-color: #0077cc;
  color: #fff;
  font-size: 16px;
  cursor: pointer;
`;

const MyComponent = () => {
  const [formData, setFormData] = useState({
    url: "",
    customAlias: "",
  });

  const handleFormChange = (event) => {
    setFormData({
      ...formData,
      [event.target.name]: event.target.value,
    });
  };

  const handleFormSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await axios.put(
        "http://localhost:8080/create",
        formData
      );
      console.log(response.data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Container>
      <Form onSubmit={handleFormSubmit}>
        <Input
          type="text"
          name="url"
          placeholder="URL"
          value={formData.url}
          onChange={handleFormChange}
        />
        <Input
          type="text"
          name="customAlias"
          placeholder="Custom Alias"
          value={formData.customAlias}
          onChange={handleFormChange}
        />
        <Button type="submit">Submit</Button>
      </Form>
    </Container>
  );
};

export default MyComponent;
