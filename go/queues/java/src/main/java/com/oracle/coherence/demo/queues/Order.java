package com.oracle.coherence.demo.queues;

import java.io.Serializable;
import java.util.Objects;

public class Order implements Serializable {
    private int orderNumber;
    private String customer;
    private String orderStatus;
    private float orderTotal;
    private long createTime;
    private long completeTime;

    public Order() {}

    public Order(int orderNumber, String customer, String orderStatus, float orderTotal, long createTime, long completeTime) {
        this.orderNumber = orderNumber;
        this.customer = customer;
        this.orderStatus = orderStatus;
        this.orderTotal = orderTotal;
        this.createTime = createTime;
        this.completeTime = completeTime;
    }

    public int getOrderNumber() {
        return orderNumber;
    }

    public void setOrderNumber(int orderNumber) {
        this.orderNumber = orderNumber;
    }

    public String getCustomer() {
        return customer;
    }

    public void setCustomer(String customer) {
        this.customer = customer;
    }

    public String getOrderStatus() {
        return orderStatus;
    }

    public void setOrderStatus(String orderStatus) {
        this.orderStatus = orderStatus;
    }

    public float getOrderTotal() {
        return orderTotal;
    }

    public void setOrderTotal(float orderTotal) {
        this.orderTotal = orderTotal;
    }

    public long getCreateTime() {
        return createTime;
    }

    public void setCreateTime(long createTime) {
        this.createTime = createTime;
    }

    public long getCompleteTime() {
        return completeTime;
    }

    public void setCompleteTime(long completeTime) {
        this.completeTime = completeTime;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        Order order = (Order) o;

        if (orderNumber != order.orderNumber) return false;
        if (Float.compare(order.orderTotal, orderTotal) != 0) return false;
        if (createTime != order.createTime) return false;
        if (completeTime != order.completeTime) return false;
        if (!Objects.equals(customer, order.customer)) return false;
        return Objects.equals(orderStatus, order.orderStatus);
    }

    @Override
    public int hashCode() {
        int result = orderNumber;
        result = 31 * result + (customer != null ? customer.hashCode() : 0);
        result = 31 * result + (orderStatus != null ? orderStatus.hashCode() : 0);
        result = 31 * result + (orderTotal != +0.0f ? Float.floatToIntBits(orderTotal) : 0);
        result = 31 * result + (int) (createTime ^ (createTime >>> 32));
        result = 31 * result + (int) (completeTime ^ (completeTime >>> 32));
        return result;
    }

    @Override
    public String toString() {
        return "Order{" +
               "orderNumber=" + orderNumber +
               ", customer='" + customer + '\'' +
               ", orderStatus='" + orderStatus + '\'' +
               ", orderTotal=" + orderTotal +
               ", createTime=" + createTime +
               ", completeTime=" + completeTime +
               '}';
    }
}
