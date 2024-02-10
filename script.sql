CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    limite INT NOT NULL,
    saldo INT NOT NULL DEFAULT 0
);


CREATE TABLE transacoes (
    id SERIAL PRIMARY KEY,
    id_cliente INT NOT NULL,
    valor INT NOT NULL,
    tipo CHAR(1) NOT NULL CHECK (tipo IN ('c', 'd')),
    descricao VARCHAR(10),
    realizada_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    limite INT NOT NULL,
    saldo INT NOT NULL,
    FOREIGN KEY (id_cliente) REFERENCES clientes(id)
);


CREATE OR REPLACE FUNCTION validar_saldo()
RETURNS TRIGGER AS $$
DECLARE
    saldo_corrente INTEGER;
    limite_cliente INTEGER;
    novo_saldo INTEGER;
BEGIN
    -- Verifiar se cliente existe
    SELECT saldo, limite INTO saldo_corrente, limite_cliente FROM clientes WHERE id = NEW.id_cliente;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'CLIENTE_NAO_ENCONTRADO';
    END IF;

    IF NEW.tipo = 'c' THEN
        novo_saldo := saldo_corrente + NEW.valor;
        
    ELSIF NEW.tipo = 'd' THEN
        novo_saldo := saldo_corrente - NEW.valor;

        -- Verificar se o saldo após a transação de débito é menor que o limite
        IF novo_saldo < (-1 * limite_cliente) THEN
            RAISE EXCEPTION 'LIMITE_EXECEDIDO';
        END IF;
    END IF;

    NEW.limite := limite_cliente;
    NEW.saldo := novo_saldo;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION atualizar_saldo()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE clientes SET saldo = NEW.saldo WHERE id = NEW.id_cliente;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER validar_saldo_trigger BEFORE INSERT ON transacoes FOR EACH ROW EXECUTE FUNCTION validar_saldo();
CREATE TRIGGER atualizar_saldo_trigger AFTER INSERT ON transacoes FOR EACH ROW EXECUTE FUNCTION atualizar_saldo();


INSERT INTO 
    clientes (nome, limite)
VALUES
    ('o barato sai caro', 1000 * 100),
    ('zan corp ltda', 800 * 100),
    ('les cruders', 10000 * 100),
    ('padaria joia de cocaia', 100000 * 100),
    ('kid mais', 5000 * 100);