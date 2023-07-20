import axios from 'axios';
import React, { useEffect, useRef, useState } from 'react';
import './App.css';

interface Todo {
  id: string
  name: string
}

const App = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [todoList, setTodoList] = useState<Todo[]>([]);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = () => {
    axios.get('/api/todo').then(res => {
      if (res.data) {
        setTodoList(res.data);
      }
    }).catch(err => {
      setTodoList([]);
      console.log(err);
    });
  }

  const addTodo = () => {
    if (inputRef.current && inputRef.current.value !== '') {
      const text = inputRef.current.value;
      axios.post('/api/todo', { name: text }).then(res => {
        if (inputRef.current) {
          inputRef.current.value = '';
          fetchData();
        }
      }).catch(err => {
        console.log(err);
      });
    }
  }

  return (
    <div className="App">
      <label>Todo</label>
      <input type="text" id="message" name="message" ref={inputRef} />
      <button onClick={addTodo}>ADD</button>
      <br/><br/><br/>
      <div className="result">
        { Array.isArray(todoList) ? todoList.map(item => <div key={item.name}>{item.name}</div>) : null}
      </div>
    </div>
  );
}

export default App;
