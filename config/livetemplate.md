| Abbreviation | Description | Template text: | variables | expression |
| ----- | ----- | ----- | ----- | ----- |
| lsv | serializeVersionUID | `    private static final long serialVersionUID = $randomCode$L;` | randomCode | groovyScript("new Random().nextLong()") |
