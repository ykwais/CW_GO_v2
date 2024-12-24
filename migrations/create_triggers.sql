
CREATE or replace TRIGGER after_user_insert
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION log_user_action();


CREATE or replace TRIGGER after_booking_insert
    AFTER INSERT ON bookings
    FOR EACH ROW
EXECUTE FUNCTION log_user_action();