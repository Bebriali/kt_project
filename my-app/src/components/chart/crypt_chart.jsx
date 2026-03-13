import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const CryptChart = () => {
    const [data, setData] = useState([]);
    const [selectedBase, setSelectedBase] = useState('SOL'); // Базовая валюта
    const [selectedQuote, setSelectedQuote] = useState('USDT'); // К чему котируется
    const [selectedExchange, setSelectedExchange] = useState('Binance');

    useEffect(() => {
        const loadData = () => {
            fetch('http://localhost:8080/stat')
                .then(res => res.json())
                .then(json => {
                    // Цены теперь уже числа (0.009219), parseFloat не обязателен,
                    // но оставим для надежности, если придет строка
                    const processedData = json.map(item => ({
                        ...item,
                        askPrice: parseFloat(item.askPrice)
                    }));
                    setData(processedData.reverse());
                })
                .catch(err => console.error("Ошибка загрузки:", err));
        };

        loadData();
        const interval = setInterval(loadData, 5000);
        return () => clearInterval(interval);
    }, []);

    // Обновленная логика фильтрации
    const filteredData = data.filter(item =>
        item.base === selectedBase &&
        item.quote === selectedQuote &&
        item.source === selectedExchange
    );

    // Динамические списки для выпадающих меню
    const bases = [...new Set(data.map(item => item.base))];
    const quotes = [...new Set(data.map(item => item.quote))];
    const exchanges = [...new Set(data.map(item => item.source))];

    return (
        <div style={{ width: '100%', height: 500, padding: '20px', backgroundColor: '#fff', borderRadius: '8px' }}>
            <h2 style={{ color: '#333' }}>Мониторинг: {selectedBase}/{selectedQuote}</h2>

            {/* Панель фильтров */}
            <div style={{ marginBottom: '20px', display: 'flex', gap: '10px' }}>
                <select value={selectedBase} onChange={(e) => setSelectedBase(e.target.value)}>
                    {bases.map(b => <option key={b} value={b}>{b}</option>)}
                </select>

                <select value={selectedQuote} onChange={(e) => setSelectedQuote(e.target.value)}>
                    {quotes.map(q => <option key={q} value={q}>{q}</option>)}
                </select>

                <select value={selectedExchange} onChange={(e) => setSelectedExchange(e.target.value)}>
                    {exchanges.map(ex => <option key={ex} value={ex}>{ex}</option>)}
                </select>
            </div>

            <ResponsiveContainer width="100%" height="80%">
                <LineChart data={filteredData}>
                    <CartesianGrid strokeDasharray="3 3" vertical={false} />
                    <XAxis
                        dataKey="timedump"
                        tickFormatter={(str) => new Date(str).toLocaleTimeString()}
                    />
                    <YAxis
                        domain={['auto', 'auto']}
                        // Увеличим точность для пар к BTC (например, 0.0012415)
                        tickFormatter={(val) => val < 1 ? val.toFixed(6) : val.toFixed(2)}
                    />
                    <Tooltip
                        labelFormatter={(label) => new Date(label).toLocaleString()}
                        contentStyle={{ borderRadius: '10px' }}
                    />
                    <Line
                        type="monotone"
                        dataKey="askPrice" // Используем новое имя поля
                        stroke="#ff0099cb"
                        strokeWidth={3}
                        dot={false}
                        activeDot={{ r: 6 }}
                        animationDuration={300}
                    />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};

export default CryptChart;
