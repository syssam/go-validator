<p>A package of validators and sanitizers for strings, structs and collections.</p>
<p>features:</p>
<ul>
  <li>Customizable Attributes.</li>
  <li>Customizable error messages.</li>
  <li>Support i18n messages</li>
</ul>
<h2>Installation</h2>
<p>Make sure that Go is installed on your computer. Type the following command in your terminal:</p>
<p>go get github.com/syssam/go-validator</p>
<h2>Usage and documentation</h2>
<h5>Examples:</h5>
<ul>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/simple">Simple</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/translations">Translations</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/gin">Gin</a></li>
</ul>
<h2>Available Validation Rules</h2>
<div class="collection-method-list">
  <p style="column-count: 3; -moz-column-count: 3; -webkit-column-count: 3; column-gap: 2em; -moz-column-gap: 2em; -webkit-column-gap: 2em;">
    <a style="display: block;">required</a>
    <a style="display: block;">requiredIf</a>
    <a style="display: block;">email</a>
    <a style="display: block;">between</a>
    <a style="display: block;">max</a>
    <a style="display: block;">min</a>
    <a style="display: block;">gt - TO DO</a>
    <a style="display: block;">gte - TO DO</a>
    <a style="display: block;">lt - TO DO</a>
    <a style="display: block;">lte - TO DO</a>
  </p>
</div>
<h4 id="rule-required">required</h4>
<p>The field under validation must be present in the input data and not empty. A field is considered "empty" if one of the following conditions are true:</p>
<div class="content-list">
  <ul>
    <li>The value is <code class="language-php"><span class="token keyword">nil</span></code>.</li>
    <li>The value is an empty string.</li>
    <li>The value is an empty array | map</li>
  </ul>
</div>
<h4 id="rule-requiredIf">requiredIf=anotherfield|value|...</h4>
<p>The field under validation must be present and not empty if the anotherfield field is equal to any value.</p>
<h4 id="rule-email">email</h4>
<p>The field under validation must be formatted as an e-mail address.</p>
<h4 id="rule-between">between=min|max</h4>
<p>The field under validation must have a size between the given min and max. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-max">max=value</h4>
<p>The field under validation must be less than or equal to a maximum value. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-min">min=value</h4>
<p>The field under validation must be greater than or equal to a minimum value. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-gt">gt=anotherfield</h4>
<p>The field under validation must be greater than the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-gte">gte=anotherfield</h4>
<p>The field under validation must be greater than or equal to the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-lt">lt=anotherfield</h4>
<p>The field under validation must be less than the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-lte">lte=anotherfield</h4>
<p>The field under validation must be less than or equal to the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h2>List of functions:</h2>
<div class="highlight highlight-source-go">
  <pre>
    IsRequiredIf(v reflect.Value, anotherfield reflect.Value, params ...string)
    IsIn(str string, params ...string) bool 
    IsEmail(str string) bool
    Between(v reflect.Value, params ...string) bool
  </pre>
</div>