const express = require('express');
const app = express();
const stripe = require('stripe')('sk_test_51KUQmKFmWCKB31rkHzO9NyfXjCvC2piQASkRfaoTHZwc3jSBOo504yqcQ7EFfjuJNO36Hjdd1cRow4ELEr90tx7D00ztlS6wm4');

app.use(express.static('public'));

app.post('/charge', async (req, res) => {
    try {
        const charge = await stripe.charges.create({
            amount: 1000, // Replace with the amount to charge in cents
            currency: 'usd',
            source: req.body.stripeToken,
            description: 'Example Charge',
        });
        res.send('Payment Successful');
    } catch (err) {
        res.status(500).send(err);
    }
});

app.listen(3000, () => console.log('Server is running on port 3000'));
