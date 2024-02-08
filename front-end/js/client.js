document.addEventListener('DOMContentLoaded', function() {
    var stripe = Stripe('pk_test_51KUQmKFmWCKB31rkR92GUPHeXv45rzjur4xWi7x4x9Sv0MgNABL7j4yRQkNM5NmzFKpyKvMXcPAt1BpBgjFUZQQU00kjLnTu01');
    var elements = stripe.elements();

    var cardElement = elements.create('card');
    cardElement.mount('#card-element');

    cardElement.on('change', function(event) {
        var displayError = document.getElementById('card-errors');
        if (event.error) {
            displayError.textContent = event.error.message;
        } else {
            displayError.textContent = '';
        }
    });

    var form = document.getElementById('payment-form');
    form.addEventListener('submit', function(event) {
        event.preventDefault();
        stripe.createToken(cardElement).then(function(result) {
            if (result.error) {
                var errorElement = document.getElementById('card-errors');
                errorElement.textContent = result.error.message;
            } else {
                stripeTokenHandler(result.token);
            }
        });
    });

    function stripeTokenHandler(token) {
        // You can send the token to your server here
        // e.g. fetch('/charge', { method: 'POST', body: token })
    }
});
